package socket

import (
	"fmt"
	"net"
	"time"
	"sync"
	"errors"
	client "github.com/balaji-balachandran/ESRG_testing/client"
)

/*
UDP Shared Socket

Client interaction with the shared socket should outwardly be the same as the single socket
Packets for clients will be buffered in queues

Client supplies its OWN buffer. Then can read from this
*/
const QUEUE_LENGTH = 10

type UDPSharedSocket struct {
	conn         	*net.UDPConn
	clientStates 	[]ClientState
	isClosed 		bool
	mu 				sync.Mutex
}

type ReadResult struct {
	packet []byte
	n int
	addr net.Addr
	err error
}

type ClientState struct {
	client   *client.UDPClient
	buffer   []byte
	callback func(srcPort int, srcIP net.IP, payload []byte) bool	// Callback that designates where the packet will be going

	queue	 chan ReadResult
	// TODO: Implement deadlines on a per connection basis
	// readDeadline time.Time
	// writeDeadline time.Time

	// Clients must be added, then loop started, then process packets
	// Look into being able to dynamically add clients (will be needed)
}

// Constructor for UDP SingleSocket
func NewUDPSharedSocket(ip net.IP, port int) (*UDPSharedSocket, error) {
	addr := &net.UDPAddr{
		IP:   ip,
		Port: port,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	sock := &UDPSharedSocket{
		conn: conn,
		isClosed: false,
		mu: sync.Mutex{},

	}

	go sock.ReadFromLoop()
	return sock, nil
}

func (sock *UDPSharedSocket) AddClient(client *client.UDPClient, buffer []byte, callback func(srcPort int, srcIP net.IP, payload []byte) bool) {
	state := ClientState{
		client:   client,
		buffer:   buffer,
		callback: callback,

		// State for queue
		queue: make(chan ReadResult, QUEUE_LENGTH),
	}
	sock.mu.Lock()
	defer sock.mu.Unlock()
	sock.clientStates = append(sock.clientStates, state)
}

// Perpetual loop to process incoming packets 
func (sock *UDPSharedSocket) ReadFromLoop(){
	for !sock.isClosed {		
		b := make([]byte, 1024)
		n, addr, err := sock.conn.ReadFromUDP(b)
		if err != nil {
			if sock.isClosed || errors.Is(err, net.ErrClosed){
				return
			}
		}
		read := ReadResult{
			packet: b,
			n: n,
			addr: addr,
			err: err,
		}
		sock.mu.Lock()
		for i, clientState := range sock.clientStates {
			// If this matches the signature expected, put
			if clientState.callback(addr.Port, addr.IP, b) {
				select{
					case clientState.queue <- read: 

					default:
						fmt.Printf("Dropping packet to client %d\n", i)
				}
				break;
			}			
		}
		sock.mu.Unlock()
	}
}

func (sock *UDPSharedSocket) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	// Check which of the clients we are trying to read from
	clientIdx := -1
	sock.mu.Lock()
	for i := range sock.clientStates {
		if len(sock.clientStates[i].buffer) > 0 && len(p) > 0 && (&(sock.clientStates[i].buffer[0]) == &(p[0])) {
			clientIdx = i
			break
		}
	}
	if clientIdx == -1 {
		sock.mu.Unlock()
		return -1, nil, nil
	}
	queue := sock.clientStates[clientIdx].queue
	sock.mu.Unlock()
	
	// Ingest from the queue 
	read_result := <- queue
	copy(p, read_result.packet)
	return read_result.n, read_result.addr, read_result.err
}

// May need to include write lock as per
// https://stackoverflow.com/questions/28575758/are-golang-net-udpconn-and-net-tcpconn-thread-safe-can-i-read-or-write-of-sing 
func (sock *UDPSharedSocket) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr.String())
	if err != nil {
		return -1, err
	}
	return sock.conn.WriteToUDP(p, udpAddr)
}

// TODO: Make it so that client close only affects per client state. Unless it is the last client and then 
func (sock *UDPSharedSocket) Close() error {
	sock.mu.Lock()
	sock.isClosed = true
	sock.mu.Unlock()
	return sock.conn.Close()
}

func (sock *UDPSharedSocket) LocalAddr() net.Addr {
	return sock.conn.LocalAddr()
}

// TODO: Implement deadline functions 
// current roadblock, no way to disambiguate clients from this framework
func (sock *UDPSharedSocket) SetDeadline(t time.Time) error {
	return sock.conn.SetDeadline(t)
}

func (sock *UDPSharedSocket) SetReadDeadline(t time.Time) error {
	return sock.conn.SetReadDeadline(t)
}

func (sock *UDPSharedSocket) SetWriteDeadline(t time.Time) error {
	return sock.conn.SetWriteDeadline(t)
}
