package main

import (
	"fmt"
	"net"
)

func main() {
	addr := net.UDPAddr{
		Port: 9000, // you can change this port
		IP:   net.ParseIP("127.0.0.1"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("UDP server listening on %s\n", addr.String())

	buf := make([]byte, 1024)
	for {
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		msg := string(buf[:n])
		go handleMessage(conn, src, msg)
	}
}

func handleMessage(conn *net.UDPConn, src *net.UDPAddr, msg string) {
	reply := fmt.Sprintf("hello %s, you said %s", src.String(), msg)
	_, err := conn.WriteToUDP([]byte(reply), src)
	fmt.Printf("reply: %v\n", reply)
	if err != nil {
		fmt.Println("Error replying:", err)
	}
}