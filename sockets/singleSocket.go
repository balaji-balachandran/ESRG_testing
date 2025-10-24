package socket

import (
	"net"
	"time"
)

type UDPSingleSocket struct {
	conn *net.UDPConn
}

// Constructor for UDP SingleSocket
func NewUDPSingleSocket(ip net.IP, port int) (*UDPSingleSocket, error) {
	addr := &net.UDPAddr{
		IP:   ip,
		Port: port,
		// Zone: zone,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &UDPSingleSocket{
		conn: conn,
	}, nil
}

func (sock *UDPSingleSocket) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return sock.conn.ReadFromUDP(p)
}

func (sock *UDPSingleSocket) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr.String())
	if err != nil {
		return -1, err
	}
	return sock.conn.WriteToUDP(p, udpAddr)
}

func (sock *UDPSingleSocket) Close() error {
	return sock.conn.Close()
}

func (sock *UDPSingleSocket) LocalAddr() net.Addr {
	return sock.conn.LocalAddr()
}

func (sock *UDPSingleSocket) SetDeadline(t time.Time) error {
	return sock.conn.SetDeadline(t)
}

func (sock *UDPSingleSocket) SetReadDeadline(t time.Time) error {
	return sock.conn.SetReadDeadline(t)
}

func (sock *UDPSingleSocket) SetWriteDeadline(t time.Time) error {
	return sock.conn.SetWriteDeadline(t)
}
