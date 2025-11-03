package tests

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

	serverAddr := net.UDPAddr{
		Port: serverPort,
		IP:   net.ParseIP(serverIP),
	}
	server_conn := server.StartServer(serverAddr, false)
	defer server_conn.Close()

	t.Run("Single Message", func(t *testing.T) {
		localIP := "127.0.0.1"
		localPort := 5005

		conn, err := socket.NewUDPSingleSocket(net.ParseIP(localIP), localPort);
		if err != nil {
			t.Fatalf("Error when creating UDPSingleSocket: %v", err)
		}
		defer conn.Close()
		// Instantiate client
		c := client.NewClient(conn, serverAddr)
		
		buf := make([]byte, 1024)

		// Send a single message
		reply, addr, err := c.SendMessage("Single message test", buf)		
		expectedReply := "[127.0.0.1:5005] Server received 'Single message test'"

		if err != nil {
			t.Fatalf("Error sending message over UDPSingleSocket: %v", err)
		}
		if reply != expectedReply{
			t.Errorf("Incorrect reply from server. Expected %s, received '%s'", expectedReply, reply)
		}
		if addr.String() != serverAddr.String() {
			t.Errorf("Incorrect server address. Expected reply from server at '%s', received from '%s'", serverAddr.String(), addr.String())
		}	
	})
}