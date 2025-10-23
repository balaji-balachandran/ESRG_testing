package receiver

import (
	"fmt"
	"net"
	"os"
	"strconv"
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

	startServer(addr)
}

func startServer( addr net.UDPAddr){
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