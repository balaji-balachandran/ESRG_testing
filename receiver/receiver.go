package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	server "github.com/balaji-balachandran/ESRG_testing/server"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Usage: go run receiver.go <IP> <Port>")
		return
	}

	port, err := strconv.Atoi(args[2])
	if err != nil {
		panic(err)
	}

	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(args[1]),
	}

	conn := server.StartServer(addr)
	defer conn.Close()

	// Block forever to keep server alive
	for {
	} 
}