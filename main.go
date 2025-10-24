package main

import (
	"fmt"
	"net"

	"github.com/balaji-balachandran/ESRG_testing/client"
	socket "github.com/balaji-balachandran/ESRG_testing/sockets"
)

func main() {
	LOCAL_IP := "127.0.0.1"
	LOCAL_PORT := 8000

	REMOTE_IP := "127.0.0.1"
	REMOTE_PORT := 9000

	remote_ip := net.ParseIP(REMOTE_IP)
	local_ip := net.ParseIP(LOCAL_IP)

	remote_addr := net.UDPAddr{
		IP:   remote_ip,
		Port: REMOTE_PORT,
	}

	sock, err := socket.NewUDPSingleSocket(local_ip, LOCAL_PORT)

	c := client.NewClient(sock, remote_addr)
	if err != nil {
		panic(err)
	}
	defer sock.Close()

	fmt.Println("Listening on", sock.LocalAddr())

	c.StdinLoop()
}
