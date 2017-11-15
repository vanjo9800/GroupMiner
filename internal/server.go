package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type client struct {
	Name       string
	IP         net.IP
	ListenPort int
	State      miningState
}

var validReqPath = regexp.MustCompile("^/(start|stop)/([0-9])$")

func makeReqHandler(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validReqPath.FindStringSubmatch(r.URL.Path)
		/*if m == nil {
			http.NotFound(w, r)
			return
		}*/
		id, err := strconv.Atoi(m[2])
		checkErr(err)
		fn(w, r, id)
	}
}

func startMinerReq(w http.ResponseWriter, r *http.Request, id int) {
	r.ParseForm()
	miningParams := mining{PoolURL: r.FormValue("poolURL"),
		Username: r.FormValue("username"), Password: r.FormValue("password"),
		Threads: converter(r.FormValue("threads")), CPUUse: converter(r.FormValue("cpuUsage"))}
	if id == 0 {
		startMiner(miningParams)
	} else {

		tcpAddr, err := net.ResolveTCPAddr("tcp", clients[id-1].IP.String()+":"+strconv.Itoa(clients[id-1].ListenPort))
		clientConn, err := net.DialTCP("tcp", nil, tcpAddr)
		defer clientConn.Close()
		for i := 0; i < 5 && err != nil; i++ {
			fmt.Println("Unable to connect to client. Trying again...")
			clientConn, err = net.DialTCP("tcp", nil, tcpAddr)
			time.Sleep(time.Second)
		}
		if err == nil {
			jsonMessage, err := json.Marshal(miningParams)
			checkErr(err)
			clientConn.Write(jsonMessage)
		} else {
			fmt.Println("Disconnected")
			//deleteFromPage(id)
		}
	}
}

func stopMinerReq(w http.ResponseWriter, r *http.Request, id int) {
	r.ParseForm()
	if id == 0 {
		stopMiner()
	} else {
		tcpAddr, err := net.ResolveTCPAddr("tcp", clients[id-1].IP.String()+":"+strconv.Itoa(clients[id-1].ListenPort))
		clientConn, err := net.DialTCP("tcp", nil, tcpAddr)
		defer clientConn.Close()
		for i := 0; i < 5 && err != nil; i++ {
			fmt.Println("Unable to connect to client. Trying again...")
			clientConn, err = net.DialTCP("tcp", nil, tcpAddr)
			time.Sleep(time.Second)
		}
		if err == nil {
			jsonMessage, err := json.Marshal(mining{PoolURL: "stop"})
			checkErr(err)
			clientConn.Write(jsonMessage)
		} else {
			fmt.Println("Disconnected")
			//deleteFromPage(id)
		}
	}
}

func statusReq(w http.ResponseWriter, r *http.Request) {
	var states []client
	mutex.Lock()
	defer mutex.Unlock()
	for id, miner := range clients {
		if id == 0 {
			mState, sysState := statusMiner()
			currentState := miningState{MiningParams: mState, SystemParams: sysState}
			states = append(states, client{Name: miner.Name, State: currentState})
		} else {
			fmt.Println(clients)
			overallState := make(chan miningState)
			go func() {
				tcpAddr, err := net.ResolveTCPAddr("tcp", clients[id].IP.String()+":"+strconv.Itoa(clients[id].ListenPort))
				checkErr(err)
				clientConn, err := net.DialTCP("tcp", nil, tcpAddr)
				checkErr(err)
				defer clientConn.Close()
				for i := 0; i < 5 && err != nil; i++ {
					fmt.Println("Unable to connect to client. Trying again...")
					clientConn, err = net.DialTCP("tcp", nil, tcpAddr)
					checkErr(err)
					time.Sleep(time.Second)
				}
				if err == nil {
					jsonMessage, err := json.Marshal(mining{PoolURL: "status"})
					checkErr(err)
					clientConn.Write(jsonMessage)
					var deviceState miningState
					data := json.NewDecoder(clientConn)
					err = data.Decode(&deviceState)
					checkErr(err)
					overallState <- deviceState

				} else {
					fmt.Println("Disconnected")
					//deleteFromPage(id)
				}
			}()
			states = append(states, client{Name: miner.Name, State: <-overallState})

		}
	}
	jsonStates, err := json.Marshal(states)
	checkErr(err)
	w.Write(jsonStates)
}

var templates = template.Must(template.ParseFiles("web/templates/index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", clients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/$")

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

var clients []client
var mutex sync.Mutex

func clientListener() {
	l, err := net.Listen("tcp", ":2000")
	checkErr(err)
	defer l.Close()
	for {
		conn, err := l.Accept()
		checkErr(err)
		data := json.NewDecoder(conn)
		var newClient clientInit
		err = data.Decode(&newClient)
		checkErr(err)
		mutex.Lock()
		clients = append(clients, client{Name: newClient.Name, ListenPort: newClient.ListenPort, IP: net.ParseIP(strings.Split(conn.RemoteAddr().String(), ":")[0])})
		mutex.Unlock()
		fmt.Println(newClient, err)
	}
}

func main() {
	clients = append(clients, client{Name: "Current device"})
	http.HandleFunc("/", makeHandler(indexHandler))
	http.HandleFunc("/start/", makeReqHandler(startMinerReq))
	http.HandleFunc("/stop/", makeReqHandler(stopMinerReq))
	http.HandleFunc("/status/", statusReq)
	http.Handle("/js/", http.FileServer(http.Dir("web/static/")))
	http.Handle("/css/", http.FileServer(http.Dir("web/static/")))

	go clientListener()
	http.ListenAndServe(":8080", nil)
}
