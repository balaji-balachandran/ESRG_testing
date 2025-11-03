package client

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
)

type UDPClient struct{
	conn net.PacketConn
	remote_addr net.UDPAddr
}

func NewClient(conn net.PacketConn, remote_addr net.UDPAddr) *UDPClient {
	return &UDPClient{
		conn: conn,
		remote_addr: remote_addr,
	}
}

func (client UDPClient) SendMessage(message string, buf []byte) (string, net.Addr, error) {

	client.conn.WriteTo([]byte(message[:]), &client.remote_addr)

	n, addr, err := client.conn.ReadFrom(buf)
	if err != nil {
		fmt.Println("Read error:", err)
		return "", nil, err
	}
	reply := string(buf[:n])

	return reply, addr, nil
}

// For debugging. Use 
func (client UDPClient) StdinLoop(buf []byte, title string){	
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("[%s] Enter message to server: ", title)
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		client.conn.WriteTo([]byte(message[:]), &client.remote_addr)

		n, addr, err := client.conn.ReadFrom(buf)
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}
		reply := string(buf[:n])
		fmt.Printf("[%s] Received from %s: %s\n", title, addr.String(), reply)
	}
}