package tests

import (
	"net"
	"testing"

	"github.com/balaji-balachandran/ESRG_testing/client"
	"github.com/balaji-balachandran/ESRG_testing/server"
	socket "github.com/balaji-balachandran/ESRG_testing/sockets"
)

func TestSharedSocketSingleConnectionSingleMessage(t *testing.T) {
	serverIP := "127.0.0.1"
	serverPort := 8000

	serverAddr := net.UDPAddr{
		Port: serverPort,
		IP:   net.ParseIP(serverIP),
	}
	server_conn := server.StartServer(serverAddr, false)
	defer server_conn.Close()


	t.Logf("Server running on %s:%d\n", serverIP, serverPort)
	t.Run("Single Message", func(t *testing.T) {
		localIP := "127.0.0.1"
		localPort := 5005
		
		expectedReply := "[127.0.0.1:5005] Server received 'Single message test'"
		conn, err := socket.NewUDPSharedSocket(net.ParseIP(localIP), localPort)
		if err != nil {
			t.Fatalf("Error when creating UDPSharedSocket: %v", err)
		}
		defer conn.Close()

		c := client.NewClient(conn, serverAddr)
		buf := make([]byte, 1024)
		callback := func(srcPort int, srcIP net.IP, p []byte) bool {
			return srcPort == serverPort && srcIP.Equal(serverAddr.IP)
		}
		conn.AddClient(c, buf, callback)

		reply, addr, err := c.SendMessage("Single message test", buf)
		if err != nil {
			t.Fatalf("Error sending message over UDPSharedSocket: %v", err)
		}
		if reply != expectedReply {
			t.Errorf("Incorrect reply from server. Expected %s, received %s", expectedReply, reply)
		}
		if addr.String() != serverAddr.String() {
			t.Errorf("Incorrect server address. Expected reply from server at %s, received from %s", serverAddr.String(), addr.String())
		}
	})
}

func TestSharedSocketSingleConnectionMultipleMessages(t *testing.T) {
	serverIP := "127.0.0.1"
	serverPort := 8000

	serverAddr := net.UDPAddr{
		Port: serverPort,
		IP:   net.ParseIP(serverIP),
	}
	server_conn := server.StartServer(serverAddr, false)
	defer server_conn.Close()

	localIP := "127.0.0.1"
	localPort := 5005

	t.Run("Multiple Messages", func(t *testing.T) {
		conn, err := socket.NewUDPSharedSocket(net.ParseIP(localIP), localPort)

		if err != nil {
			t.Fatalf("Error when creating UDPSharedSocket: %v", err)
		}
		defer conn.Close()

		c := client.NewClient(conn, serverAddr)
		buf := make([]byte, 1024)
		callback := func(srcPort int, srcIP net.IP, p []byte) bool {
			return srcPort == serverPort && srcIP.Equal(serverAddr.IP)
		}
		conn.AddClient(c, buf, callback)

		expectedReply := "[127.0.0.1:5005] Server received 'Multiple message test [1/2]'"
		reply, addr, err := c.SendMessage("Multiple message test [1/2]", buf)
		if err != nil {
			t.Fatalf("Error sending message over UDPSingleSocket: %v", err)
		}
		if reply != expectedReply {
			t.Errorf("Incorrect reply from server. Expected %s, received %s", expectedReply, reply)
		}
		if addr.String() != serverAddr.String() {
			t.Errorf("Incorrect server address. Expected reply from server at %s, received from %s", serverAddr.String(), addr.String())
		}

		expectedReply2 := "[127.0.0.1:5005] Server received 'Multiple message test [2/2]'"

		reply, addr, err = c.SendMessage("Multiple message test [2/2]", buf)
		if err != nil {
			t.Fatalf("Error sending message over UDPSharedSocket: %v", err)
		}
		if reply != expectedReply2 {
			t.Errorf("Incorrect reply from server. Expected %s, received %s", expectedReply2, reply)
		}
		if addr.String() != serverAddr.String() {
			t.Errorf("Incorrect server address. Expected reply from server at %s, received from %s", serverAddr.String(), addr.String())
		}
	})
}

func TestSharedSocketMultipleConnectionsSingleMessage(t *testing.T) {
	server1IP := "127.0.0.1"
	server1Port := 8000

	server1Addr := net.UDPAddr{
		Port: server1Port,
		IP:   net.ParseIP(server1IP),
	}
	server1_conn := server.StartServer(server1Addr, false)
	defer server1_conn.Close()

	server2IP := "127.0.0.1"
	server2Port := 9000

	server2Addr := net.UDPAddr{
		Port: server2Port,
		IP:   net.ParseIP(server2IP),
	}
	server2_conn := server.StartServer(server2Addr, false)
	defer server2_conn.Close()

	localIP := "127.0.0.1"
	localPort := 5005

	t.Run("Multiple Connections", func(t *testing.T) {
		conn, err := socket.NewUDPSharedSocket(net.ParseIP(localIP), localPort)

		if err != nil {
			t.Fatalf("Error when creating UDPSharedSocket: %v", err)
		}
		defer conn.Close()

		c1 := client.NewClient(conn, server1Addr)
		buf := make([]byte, 1024)
		callback := func(srcPort int, srcIP net.IP, p []byte) bool {
			return srcPort == server1Port && srcIP.Equal(server1Addr.IP)
		}
		conn.AddClient(c1, buf, callback)

		c2 := client.NewClient(conn, server2Addr)
		buf2 := make([]byte, 1024)
		callback2 := func(srcPort int, srcIP net.IP, p []byte) bool {
			return srcPort == server2Port && srcIP.Equal(server2Addr.IP)
		}
		conn.AddClient(c2, buf2, callback2)

		expectedReply := "[127.0.0.1:5005] Server received 'Single message test, sharing socket [client 1]'"
		t.Log("Client 1 attempting send")
		reply, addr, err := c1.SendMessage("Single message test, sharing socket [client 1]", buf)
		t.Log("Client 1 received response")
		if err != nil {
			t.Fatalf("Error sending message over UDPSharedSocket: %v", err)
		}
		if reply != expectedReply {
			t.Errorf("Incorrect reply from server. Expected %s, received %s", expectedReply, reply)
		}
		if addr.String() != server1Addr.String() {
			t.Errorf("Incorrect server address. Expected reply from server at %s, received from %s", server1Addr.String(), addr.String())
		}

		expectedReply2 := "[127.0.0.1:5005] Server received 'Single message test, sharing socket [client 2]'"

		t.Log("Client 2 attempting send")
		reply, addr, err = c2.SendMessage("Single message test, sharing socket [client 2]", buf2)
		t.Log("Client 2 received response")
		if err != nil {
			t.Fatalf("Error sending message over UDPSharedSocket: %v", err)
		}
		if reply != expectedReply2 {
			t.Errorf("Incorrect reply from server. Expected %s, received %s", expectedReply2, reply)
		}
		if addr.String() != server2Addr.String() {
			t.Errorf("Incorrect server address. Expected reply from server at %s, received from %s", server2Addr.String(), addr.String())
		}
	})
}