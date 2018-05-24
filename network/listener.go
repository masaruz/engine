package network

import (
	"fmt"
	"log"
	"net"

	"github.com/masaruz/engine-lib/core"
)

// ListenAndServe requests
func ListenAndServe(port int, game core.Game) error {
	// Lets prepare a address at any address at port 10001
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	// Now listen at selected port
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		return err
	}

	defer ServerConn.Close()
	// Get interface addresses in container
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}
	// Logging addresses
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Printf("Server is running at %s:%d\n", ipnet.IP.String(), port)
			}
		}
	}
	// Waiting for requests
	buf := make([]byte, 1024)
	session := CreateSession(ServerConn)
	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		msg := buf[0:n]

		// Check if addr is already joined session
		// Register client who send packet
		session.Join(addr)

		// TODO ...
		// If msg is acknowledge message
		if false {
			session.ACK("packetID")
		}
		fmt.Println("Received", string(msg), "from", addr, "length", n)
		// Send update to game
		// and receive acknowledge
		if err := game.Update(msg, func(ack string) {
			// do something ...
		}); err != nil {
			log.Println("Error: ", err)
		}
		// Send back to client
		session.Send(game.GetState())
	}
}
