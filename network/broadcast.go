package network

import (
	"net"
	"sync"
)

// Broadcast handle client's addresses
type Broadcast struct {
	c    map[string]*net.UDPAddr
	m    sync.Mutex
	Conn *net.UDPConn
}

// Join add client to map
func (b *Broadcast) Join(addr *net.UDPAddr) error {
	if b.c == nil {
		b.c = make(map[string]*net.UDPAddr)
	}
	if b.c[addr.String()] == nil {
		b.m.Lock()
		b.c[addr.String()] = addr
		b.m.Unlock()
	}
	return nil
}

// Leave remove client from map
func (b *Broadcast) Leave(addr *net.UDPAddr) error {
	if b.c == nil {
		b.c = make(map[string]*net.UDPAddr)
	}
	if b.c[addr.String()] != nil {
		b.m.Lock()
		delete(b.c, addr.String())
		b.m.Unlock()
	}
	return nil
}

// Send all message to client
func (b *Broadcast) Send(payload []byte) error {
	for addr := range b.c {
		b.Conn.WriteToUDP(payload, b.c[addr])
	}
	return nil
}
