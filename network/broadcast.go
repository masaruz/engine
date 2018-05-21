package network

import (
	"log"
	"net"
	"sync"
)

// Session handle client's addresses
type Session struct {
	c    map[string]*net.UDPAddr
	m    sync.Mutex
	Conn *net.UDPConn
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

// Send all message to client
func (s *Session) Send(payload []byte) error {
	for addr := range s.c {
		log.Println("Send to", s.c[addr])
		s.Conn.WriteToUDP(payload, s.c[addr])
	}
	return nil
}
