package main

import (
	"net"
	"testing"

	"github.com/balaji-balachandran/ESRG_testing/client"
	"github.com/balaji-balachandran/ESRG_testing/server"
	"github.com/balaji-balachandran/ESRG_testing/sockets"
)

func TestSingleSocketConnection(t *testing.T){
	serverIP := "127.0.0.1"
	serverPort := 8000

	server_addr := net.UDPAddr{
		Port: serverPort,
		IP:   net.ParseIP(serverIP),
	}
	server_conn := server.StartServer(server_addr)
	defer server_conn.Close()

	local_ip := "127.0.0.1"
	local_port := 5005
	expected_reply := "[127.0.0.1:5005] Server received 'Single message test'"
	
	t.Run("Single Message", func(t *testing.T) {
		conn, err := socket.NewUDPSingleSocket(net.ParseIP(local_ip), local_port);
		if err != nil {
			t.Fatalf("Error when creating UDPSingleSocket: %v", err)
		}
		c := client.NewClient(conn, server_addr)
		
		reply, addr, err := c.SendMessage("Single message test")
		if err != nil {
			t.Fatalf("Error sending message over UDPSingleSocket: %v", err)
		}
		if reply != expected_reply{
			t.Errorf("Incorrect reply from server. Expected %s, received %s", expected_reply, reply)
		}
		if addr.String() != server_addr.String() {
			t.Errorf("Incorrect server address. Expected reply from server at %s, received from %s", server_addr.String(), addr.String())
		}	
	})
}