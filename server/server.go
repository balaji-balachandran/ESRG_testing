package server

import (
	"net"
	"fmt"
)

func StartServer(addr net.UDPAddr) *net.UDPConn {
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("UDP server listening on %s\n", addr.String())

	go func() {
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
	}()

	return conn
}

func handleMessage(conn *net.UDPConn, src *net.UDPAddr, msg string) {
	reply := fmt.Sprintf("[%s] Server received '%s'", src.String(), msg)
	_, err := conn.WriteToUDP([]byte(reply), src)
	fmt.Printf("reply: %v\n", reply)
	if err != nil {
		fmt.Println("Error replying:", err)
	}
}
