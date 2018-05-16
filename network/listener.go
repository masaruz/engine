package network

import (
	"fmt"
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

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Printf("Server is running at %s:%d\n", ipnet.IP.String(), port)
			}
		}
	}
	// Waiting for requests
	buf := make([]byte, 1024)
	b := &Broadcast{Conn: ServerConn}
	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		b.Join(addr)
		msg := buf[0:n]
		fmt.Println("Received", string(msg), "from", addr, "length", n)
		if err := game.Update(msg); err != nil {
			panic(err)
		}
		// Send back to client
		b.Send(msg)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
