package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/shirou/gopsutil/cpu"
)

type mining struct {
	PoolURL  string
	Username string
	Password string
	Threads  int
	CPUUse   int
}

type system struct {
	CurrentState bool
	ProcessCPU   float64
	SystemCPU    float64
	Temp         float64
	ProcessRAM   float64
	SystemRAM    float64
}

type miningState struct {
	MiningParams mining
	SystemParams system
}

type client struct {
	Name       string
	Type       string
	IP         net.IP
	ListenPort int
	State      miningState
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func converter(s string) (i int) {
	val, err := strconv.Atoi(s)
	checkErr(err)
	return val
}

func startMiner(m mining) {
	miner := exec.CommandContext(context.Background(), "/bin/bash",
		"scripts/m-minerd.sh", m.PoolURL, m.Username, m.Password, strconv.Itoa(m.Threads), strconv.Itoa(m.CPUUse))
	miner.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := miner.Start()
	checkErr(err)
	fmt.Printf("Started miner with PID:%d\n", miner.Process.Pid)
	err = ioutil.WriteFile("miner.pid", []byte(strconv.Itoa(miner.Process.Pid)), 0777)
	checkErr(err)
	err = miner.Process.Release()
	checkErr(err)
}

func stopMiner() {
	minerPid, err := ioutil.ReadFile("miner.pid")
	os.Remove("miner.pid")
	checkErr(err)
	pid, err := strconv.Atoi(string(minerPid))
	checkErr(err)
	fmt.Printf("Killed miner with PID:%d\n", pid)
	pgid, err := syscall.Getpgid(pid)
	checkErr(err)
	syscall.Kill(-pgid, syscall.SIGKILL)
}

func statusMiner() (m mining, s system) {
	if _, err := os.Stat("miner.pid"); os.IsNotExist(err) {
		return mining{}, system{CurrentState: false}
	}
	minerPid, err := ioutil.ReadFile("miner.pid")
	checkErr(err)
	pid, err := strconv.Atoi(string(minerPid))
	checkErr(err)
	argsBytes, err := ioutil.ReadFile("/proc/" + strconv.Itoa(pid) + "/cmdline")
	zeroByte := make([]byte, 1)
	zeroByte[0] = 0
	args := bytes.Split(argsBytes, zeroByte)
	CPUUsage, err := cpu.Percent(0, false)
	checkErr(err)
	return mining{PoolURL: string(args[2]), Username: string(args[3]), Threads: converter(string(args[5])),
		CPUUse: converter(string(args[6]))}, system{CurrentState: true, SystemCPU: float64(int(CPUUsage[0]*10)) / 10}
}
