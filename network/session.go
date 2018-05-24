package network

import (
	"net"
	"os"
	"sync"
	"time"

	"github.com/masaruz/engine-lib/common"
)

const (
	maxRetry = 3
)

// Session handle client's addresses
type Session struct {
	ack  map[string]chan bool
	c    map[string]*net.UDPAddr
	m    sync.Mutex
	Conn *net.UDPConn
}

// CreateSession with udp and ack
func CreateSession(ServerConn *net.UDPConn) *Session {
	return &Session{
		Conn: ServerConn,
		ack:  make(map[string]chan bool),
	}
}

// Join add client to map
func (s *Session) Join(addr *net.UDPAddr) error {
	if s.c == nil {
		s.c = make(map[string]*net.UDPAddr)
	}
	if s.c[addr.String()] == nil {
		s.m.Lock()
		s.c[addr.String()] = addr
		s.m.Unlock()
	}
	return nil
}

// Leave remove client from map
func (s *Session) Leave(addr *net.UDPAddr) error {
	if s.c == nil {
		s.c = make(map[string]*net.UDPAddr)
	}
	if s.c[addr.String()] != nil {
		s.m.Lock()
		delete(s.c, addr.String())
		s.m.Unlock()
	}
	return nil
}

// Send all clients the message
func (s *Session) Send(payload []byte) error {
	for addr := range s.c {
		common.Print("Send to", s.c[addr])
		// TODO ...
		// If this broadcast need to be delivered
		if false {
			s.create("packetID")
			go s.wait("packetID", addr, payload, maxRetry)
		}
		common.Print("Write to", addr)
		s.write(payload, addr)
	}
	return nil
}

// ACK send to channel to confirm that send successfully
func (s *Session) ACK(packetID string) {
	if s.ack[packetID] != nil {
		s.ack[packetID] <- true
	}
}

// Create channel each packet which need to be delivered
func (s *Session) create(packetID string) chan bool {
	if s.ack[packetID] == nil {
		s.ack[packetID] = make(chan bool)
	}
	return s.ack[packetID]
}

// Wait for ack to handle reliable messages
func (s *Session) wait(packetID string, addr string, payload []byte, retry int) {
	// Exceed limit retry then delete from ack list
	if retry == 0 {
		delete(s.ack, packetID)
		return
	}
	select {
	case <-s.ack[packetID]:
		// Successfully sent
		delete(s.ack, packetID)
	case <-time.After(time.Second):
		// Timeout and try to resend the package
		s.write(payload, addr)
		s.wait(packetID, addr, payload, retry-1)
	}
}

// Write message via UDP
func (s *Session) write(payload []byte, addr string) {
	if os.Getenv("stage") == "dev" {
		return
	}
	s.Conn.WriteToUDP(payload, s.c[addr])
}
