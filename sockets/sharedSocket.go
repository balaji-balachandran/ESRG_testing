package socket

import (
	"net/netip"
	"time"
	"sync"
	"errors"
	"net"
	"sync/atomic"
	"log/slog"
	"fmt"
)

/*
UDP Shared Socket

Client interaction with the shared socket should outwardly be the same as the single socket
Packets for clients will be buffered in queues

Client supplies its OWN buffer. Then can read from this
*/
const QUEUE_LENGTH = 10

type FlowKey struct {
	SrcIP, DstIP   netip.Addr
	SrcPort, DstPort uint16
}

type SharedSocketDialer struct {
	conn    net.PacketConn
	conns 	[]*SharedSocketConn
	mu 		sync.RWMutex
	network string
	closed  atomic.Bool
}

type SharedSocketConn struct {
	cb func (network string, srcIP netip.Addr, srcPort uint16, actualPacket []byte) bool
	flowKey FlowKey // the unique flow key for this connection, should be immutable. Used to clean up the parent's maps on close
	parent *SharedSocketDialer // reference to the parent shared socket
	inCh   chan ReadResult // channel for incoming packets, buffered
	// writes can be done directly to the parent socket, go's net.PacketConn is thread safe
	closed atomic.Bool // to prevent double closes.
}

type ReadResult struct {
	packet []byte
	n int
	addr net.Addr
	err error
}

// Constructor for UDP SingleSocket
// This should be an interface that we have udp version, icmp version, etc.
func NewUDPSharedSocketDialer(localAddr *net.UDPAddr) (*SharedSocketDialer, error) {

	// setup the common state here, init locks/maps, get the shared socket open, etc.
	// Listen > Dial, Dial only supports a single remote address (which means we cannot use
	// it for a shared socket)
	conn, err := net.ListenUDP(localAddr.Network(), localAddr)

	if err != nil {
		slog.Error("unable to listen on port")
		return nil, err
	}
	
	return &SharedSocketDialer{
		conn: conn,
		conns: []*SharedSocketConn{},
		mu: sync.RWMutex{},
		network: localAddr.Network(),
	}, nil
}

// Adds a new connection to the sharedSocketDialer
func (s *SharedSocketDialer) Dial(network, address string, callback func (network string, srcIP netip.Addr, srcPort uint16, actualPacket []byte) bool, flowkey FlowKey) (*SharedSocketConn, error) {

	// If proposed network is not same as the Dialer's network, exit and
	// return error about the mismatch
	if network != s.network {
		// TODO: Make more descriptive error
		return nil, &net.AddrError{}
	}

	// Expose a client connection to the client
	client_conn := &SharedSocketConn{
		cb : callback,
		flowKey: flowkey,
		parent: s,
		inCh: make(chan ReadResult, QUEUE_LENGTH),
		closed: atomic.Bool{}, 			// Zero value is false
	}
	
	// Add this client connection
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conns = append(s.conns, client_conn)

	return client_conn, nil
}

// Perpetual loop to process incoming packets 
func (sock *SharedSocketDialer) ReadFromLoop(){
	for {		
		b := make([]byte, 1024)
		n, addr, err := sock.conn.ReadFrom(b)
		if err != nil {
			if errors.Is(err, net.ErrClosed){
				return
			}
		}
		read := ReadResult{
			packet: b,
			n: n,
			addr: addr,
			err: err,
		}

		var srcIP netip.Addr;
		var srcPort uint16;

		// Conditional based on type of net.addr
		switch a := addr.(type) {
		case *net.UDPAddr:
			ip, ok := netip.AddrFromSlice(a.IP)
			if ok {
				srcIP = ip
				srcPort = netip.AddrPortFrom(ip, uint16(a.Port)).Port()
			}
		case *net.IPAddr:
			ip, ok := netip.AddrFromSlice(a.IP)
			if ok {
				srcIP = ip
				srcPort = 0 // IP addrs have no ports, set to 0
			}
		}


		sock.mu.RLock()
		// Check the client connections to see which call back matches and then forward to that connection
		for _, conn := range(sock.conns){
			if conn.cb(sock.network, srcIP, srcPort, b){
				select {
					case conn.inCh <- read: 
					default:
						fmt.Printf("Dropping packet to client %v\n", conn.flowKey)
				}
				break
			}
		}
		sock.mu.RUnlock()
	}
}

func (c *SharedSocketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	readResult := <- c.inCh
	copy(p, readResult.packet)
	return readResult.n, readResult.addr, readResult.err
}

// May need to include write lock as per
// https://stackoverflow.com/questions/28575758/are-golang-net-udpconn-and-net-tcpconn-thread-safe-can-i-read-or-write-of-sing 
func (c *SharedSocketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return c.parent.conn.WriteTo(p, addr)
}

// TODO: Make it so that client close only affects per client state. Unless it is the last client and then 
func (c *SharedSocketConn) Close() error {
	// Connection already closed, do nothing (or TODO: return error)
	if c.closed.Load() {
		return nil
	}
	c.closed.Store(true)
	
	c.parent.mu.Lock()
	defer c.parent.mu.Unlock()
	
	removeIndex := -1
	for i, conn := range(c.parent.conns){
		if conn == c {
			removeIndex = i
		}
	}
	c.parent.conns[removeIndex] = c.parent.conns[len(c.parent.conns) - 1]
	c.parent.conns = c.parent.conns[:len(c.parent.conns) - 1]

	
	if(len(c.parent.conns) > 0) { return nil }
	
	// If we have no remaining client connections, truly close the shared socket connection
	return c.parent.conn.Close()
}

func (c *SharedSocketConn) LocalAddr() net.Addr {
	return c.parent.conn.LocalAddr()
}

// TODO: Implement deadline functions 
// current roadblock, no way to disambiguate clients from this framework
func (c *SharedSocketConn) SetDeadline(t time.Time) error {
	return c.parent.conn.SetDeadline(t)
}

func (c *SharedSocketConn) SetReadDeadline(t time.Time) error {
	return c.parent.conn.SetReadDeadline(t)
}

func (c *SharedSocketConn) SetWriteDeadline(t time.Time) error {
	return c.parent.conn.SetWriteDeadline(t)
}
