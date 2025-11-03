package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"github.com/balaji-balachandran/ESRG_testing/client"
	socket "github.com/balaji-balachandran/ESRG_testing/sockets"
)

func main() {
	args := os.Args
	if(len(args) < 5) {
		fmt.Println("Usage: go run main.go <Local_IP> <Local_Port> <Remote_IP> <Remote_Port> (<Remote_IP> <Remote_Port>...)")
		os.Exit(1)
	}
	LOCAL_IP := args[1]
	LOCAL_PORT, err := strconv.Atoi(args[2])

	REMOTE_IP := args[3]
	REMOTE_PORT, err  := strconv.Atoi(args[4])

	REMOTE_IP2 := args[5]
	REMOTE_PORT2, err := strconv.Atoi(args[6])

	local_ip := net.ParseIP(LOCAL_IP)
	remote_ip := net.ParseIP(REMOTE_IP)
	remote_ip2 := net.ParseIP(REMOTE_IP2)

	remote_addr := net.UDPAddr{
		IP:   remote_ip,
		Port: REMOTE_PORT,
	}

	remote_addr2 := net.UDPAddr{
		IP: remote_ip2,
		Port: REMOTE_PORT2,
	}

	sock, err := socket.NewUDPSharedSocket(local_ip, LOCAL_PORT)
	defer sock.Close()

	if err != nil {
		panic(err)
	}
	// Make a new client
	c := client.NewClient(sock, remote_addr)
	c2 := client.NewClient(sock, remote_addr2)
	// Add client to shared socket buffer
	buf := make([]byte, 1024)
	callback := func(srcPort int, srcIP net.IP, payload []byte) bool{
		return (srcPort == REMOTE_PORT) && srcIP.Equal(remote_ip)
	}
	sock.AddClient(c, buf, callback)


	buf2 := make([]byte, 1024)
	callback2 := func(srcPort int, srcIP net.IP, payload []byte) bool {
		return (srcPort == REMOTE_PORT2) && srcIP.Equal(remote_ip2)
	}
	sock.AddClient(c2, buf2, callback2)
	fmt.Println("Listening on", sock.LocalAddr())

	go c.StdinLoop(buf, "client 1")
	go c2.StdinLoop(buf2, "client 2")
	for {

	}
}
