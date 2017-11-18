#!/bin/bash

mkdir -p log/
touch log/setup.log
echo "If any errors occur during the setup, report them via the log located at log/setup.log"
echo "Downloading Go..."
echo "Please select an operating system:"
echo "1) Ubuntu, or other DEB-based"
echo "2) Arch Linux"
echo "3) Other Linux distribution (you will need to install yourself)"
echo -n "Select (1/2/None):" && read n
case $n in
    1) sudo apt-get install golang;;
    2) sudo pacman -S go;;
esac
export GOPATH=$HOME/go
echo "Downloading the libraries..."
go get github.com/shirou/gopsutil/cpu >> log/setup.log
go get github.com/gorilla/websocket >> log/setup.log
echo "Building Go files..."
mkdir -p bin
go build internal/server.go internal/miner.go >> log/setup.log
go build internal/client.go internal/miner.go >> log/setup.log
mv server bin/
mv client bin/
echo "Successfully finished configuring the Go environment!"
echo "Now you can start the service via bin/server, or bin/client depending on the working mode."
