package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type serverConfig struct {
	WebPort    int
	ClientPort int
}

func readConf() (conf serverConfig) {
	configFile, err := ioutil.ReadFile("config/server.conf")
	checkErr(err)
	configs := bytes.Split(configFile, []byte("\n"))
	return serverConfig{WebPort: converter(strings.Split(string(configs[1]), "=")[1]), ClientPort: converter(strings.Split(string(configs[2]), "=")[1])}
}

var webClients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var validReqPath = regexp.MustCompile("^/(start|stop)/([0-9])$")

func makeReqHandler(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validReqPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
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
		tcpAddr, err := net.ResolveTCPAddr("tcp", clients[id].IP.String()+":"+strconv.Itoa(clients[id].ListenPort))
		checkErr(err)
		clientConn, err := net.DialTCP("tcp", nil, tcpAddr)
		for i := 0; i < 5 && err != nil; i++ {
			fmt.Println("Unable to connect to client " + clients[id].Name + ". Trying again...")
			clientConn, err = net.DialTCP("tcp", nil, tcpAddr)
			time.Sleep(time.Second)
		}
		if err == nil {
			jsonMessage, err := json.Marshal(miningParams)
			checkErr(err)
			clientConn.Write(jsonMessage)
			clientConn.Close()
		} else {
			fmt.Println("Disconnected client " + clients[id].Name + " from " + clients[id].IP.String())
			clients = append(clients[:id], clients[id+1:]...)
		}
	}
}

func stopMinerReq(w http.ResponseWriter, r *http.Request, id int) {
	r.ParseForm()
	if id == 0 {
		stopMiner()
	} else {
		tcpAddr, err := net.ResolveTCPAddr("tcp", clients[id].IP.String()+":"+strconv.Itoa(clients[id].ListenPort))
		checkErr(err)
		clientConn, err := net.DialTCP("tcp", nil, tcpAddr)
		for i := 0; i < 5 && err != nil; i++ {
			fmt.Println("Unable to connect to client " + clients[id].Name + ". Trying again...")
			clientConn, err = net.DialTCP("tcp", nil, tcpAddr)
			time.Sleep(time.Second)
		}
		if err == nil {
			jsonMessage, err := json.Marshal(mining{PoolURL: "stop"})
			checkErr(err)
			clientConn.Write(jsonMessage)
			clientConn.Close()
		} else {
			fmt.Println("Disconnected client " + clients[id].Name + " from " + clients[id].IP.String())
			clients = append(clients[:id], clients[id+1:]...)
		}
	}
}

func sendStatus() {
	var states []client
	mutex.Lock()
	defer mutex.Unlock()
	for id, miner := range clients {
		if id == 0 {
			mState, sysState := statusMiner()
			currentState := miningState{MiningParams: mState, SystemParams: sysState}
			states = append(states, client{Name: miner.Name, IP: miner.IP, State: currentState})
		} else {
			overallState := make(chan miningState)
			success := make(chan bool)
			go func() {
				tcpAddr, err := net.ResolveTCPAddr("tcp", miner.IP.String()+":"+strconv.Itoa(miner.ListenPort))
				checkErr(err)
				clientConn, err := net.DialTCP("tcp", nil, tcpAddr)
				for i := 0; i < 5 && err != nil; i++ {
					fmt.Println("Unable to connect to client " + clients[id].Name + ". Trying again...")
					clientConn, err = net.DialTCP("tcp", nil, tcpAddr)
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
					success <- true
					overallState <- deviceState
					clientConn.Close()
				} else {
					fmt.Println("Disconnected client " + clients[id].Name + " from " + clients[id].IP.String())
					clients = append(clients[:id], clients[id+1:]...)
					success <- false
				}
			}()
			if <-success {
				states = append(states, client{Name: miner.Name, Type: miner.Type, IP: miner.IP, State: <-overallState})
			}
		}
	}
	for webClient := range webClients {
		err := webClient.WriteJSON(states)
		if err != nil {
			webClient.Close()
			delete(webClients, webClient)
		}
	}
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

func clientListener(conf serverConfig) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(conf.ClientPort))
	checkErr(err)
	defer l.Close()
	for {
		conn, err := l.Accept()
		checkErr(err)
		data := json.NewDecoder(conn)
		var newClient client
		err = data.Decode(&newClient)
		checkErr(err)
		mutex.Lock()
		remoteIP := net.ParseIP(strings.Split(conn.RemoteAddr().String(), ":")[0])
		inside := false
		for _, existingClient := range clients {
			if existingClient.Name == newClient.Name && existingClient.IP.String() == remoteIP.String() {
				inside = true
			}
		}
		if !inside {
			fmt.Println("Added client " + newClient.Name + " at " + remoteIP.String() + ":" + strconv.Itoa(newClient.ListenPort))
			clients = append(clients, client{Name: newClient.Name, Type: newClient.Type, ListenPort: newClient.ListenPort, IP: remoteIP})
		}
		mutex.Unlock()
	}
}

func handleWebsockets(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	checkErr(err)
	//defer ws.Close()
	webClients[ws] = true
}

func main() {
	currentConfig := readConf()
	fmt.Println("Starter GroupMiner server on http://localhost:" + strconv.Itoa(currentConfig.WebPort) + "/")
	clients = append(clients, client{Name: "Current device", IP: net.ParseIP("127.0.0.1")})
	http.HandleFunc("/", makeHandler(indexHandler))
	http.HandleFunc("/start/", makeReqHandler(startMinerReq))
	http.HandleFunc("/stop/", makeReqHandler(stopMinerReq))
	http.Handle("/js/", http.FileServer(http.Dir("web/static/")))
	http.Handle("/css/", http.FileServer(http.Dir("web/static/")))
	http.Handle("/images/", http.FileServer(http.Dir("web/static/")))
	http.HandleFunc("/ws", handleWebsockets)

	go func() {
		for {
			sendStatus()
			time.Sleep(time.Second)
		}
	}()
	go clientListener(currentConfig)
	http.ListenAndServe(":"+strconv.Itoa(currentConfig.WebPort), nil)
}
