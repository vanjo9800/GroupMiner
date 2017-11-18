package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"
)

type clientConfig struct {
	Name       string
	Type       string
	Server     string
	ServerPort int
	ListenPort int
}

func readConf() (conf clientConfig) {
	configFile, err := ioutil.ReadFile("config/client.conf")
	checkErr(err)
	configs := bytes.Split(configFile, []byte("\n"))
	return clientConfig{Name: strings.Split(string(configs[1]), "=")[1], Type: strings.Split(string(configs[2]), "=")[1], Server: strings.Split(string(configs[3]), "=")[1],
		ServerPort: converter(strings.Split(string(configs[4]), "=")[1]), ListenPort: converter(strings.Split(string(configs[5]), "=")[1])}
}

var connected bool

func connectToServer(conf clientConfig) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", conf.Server+":"+strconv.Itoa(conf.ServerPort))
	serverConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err == nil {
		if !connected {
			fmt.Println("Connected to server at http://" + conf.Server + ":" + strconv.Itoa(conf.ServerPort) + "/")
			connected = true
		}
		defer serverConn.Close()
		jsonMessage, err := json.Marshal(client{Name: conf.Name, Type: conf.Type, ListenPort: conf.ListenPort})
		checkErr(err)
		serverConn.Write(jsonMessage)
	} else {
		connected = false
		fmt.Println("Unable to connect to server " + tcpAddr.String() + ". Trying again...")
	}
}

func startListener(conf clientConfig) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(conf.ListenPort))
	checkErr(err)
	defer l.Close()
	for {
		conn, err := l.Accept()
		checkErr(err)
		data := json.NewDecoder(conn)
		var miningDetails mining
		err = data.Decode(&miningDetails)
		if miningDetails.PoolURL == "status" {
			mState, sysState := statusMiner()
			overallState := miningState{MiningParams: mState, SystemParams: sysState}
			jsonMessage, err := json.Marshal(overallState)
			checkErr(err)
			conn.Write(jsonMessage)
		} else {
			if miningDetails.PoolURL == "stop" {
				stopMiner()
			} else {
				startMiner(miningDetails)
			}
		}
	}
}

func main() {
	currentConfig := readConf()
	fmt.Println("Starter GroupMiner client")
	go func() {
		for {
			connectToServer(currentConfig)
			time.Sleep(time.Second)
		}
	}()
	startListener(currentConfig)
}
