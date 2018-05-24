package network

import (
	"net"
	"testing"
	"time"
)

var (
	ip1 = &net.UDPAddr{IP: []byte("1")}
	ip2 = &net.UDPAddr{IP: []byte("2")}
)

func TestSessionJoin(t *testing.T) {
	session := CreateSession(nil)
	if err := session.Join(ip1); len(session.c) != 1 || err != nil {
		t.Error()
	}
	if err := session.Join(ip2); len(session.c) != 2 || err != nil {
		t.Error()
	}
}
func TestSessionLeave(t *testing.T) {
	session := CreateSession(nil)
	if err := session.Leave(ip1); err != nil {
		t.Error()
	}
	if err := session.Leave(ip2); err != nil {
		t.Error()
	}
	if err := session.Join(ip1); len(session.c) != 1 || err != nil {
		t.Error()
	}
	if err := session.Join(ip2); len(session.c) != 2 || err != nil {
		t.Error()
	}
	if err := session.Leave(ip1); len(session.c) != 1 || err != nil {
		t.Error()
	}
	if err := session.Leave(ip2); len(session.c) != 0 || err != nil {
		t.Error()
	}
}
func TestSessionWait(t *testing.T) {
	session := CreateSession(nil)
	packetID := "id"
	session.create(packetID)
	go session.wait(packetID, "localhost", []byte{}, 3)
	if session.ack[packetID] == nil {
		t.Error()
	}
	// Assume that client ack on time
	session.ack[packetID] <- true
	time.Sleep(time.Millisecond)
	if session.ack[packetID] != nil {
		t.Error()
	}
	session.create(packetID)
	go session.wait(packetID, "localhost", []byte{}, 1)
	// Assume that client no response until timeout
	time.Sleep(time.Millisecond * 1100)
	if session.ack[packetID] != nil {
		t.Error()
	}
}
