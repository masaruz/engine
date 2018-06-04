package network

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/masaruz/engine-lib/common"
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
	fmt.Printf("Server is running at port: %d\n", port)
	// Waiting for requests
	session := CreateSession(ServerConn)
	// Loop listening and reading from UDP
	go listen(ServerConn, session, game)
	// Loop sending to client
	send(session, game)
	return nil
}

func send(session *Session, game core.Game) {
	sleep := 200 // Default cooldown
	// If send cooldown is set
	if conv, err := strconv.Atoi(os.Getenv("SEND_COOLDOWN")); err == nil {
		sleep = conv
	}
	for {
		// Send back to client
		go session.Send(game.GetState())
		time.Sleep(time.Millisecond * time.Duration(sleep))
	}
}

func listen(ServerConn *net.UDPConn, session *Session, game core.Game) {
	buf := make([]byte, 1024)
	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error: ", err)
		}
		msg := buf[0:n]
		// Check if addr is already joined session
		// Register client who send packet
		session.Join(addr)
		// TODO ...
		// If msg is acknowledge message
		if false {
			go session.ACK("packetID")
			common.Print("Received acknowledge", string(msg), "from", addr)
		} else {
			common.Print("Received", string(msg), "from", addr)
			// Send update to game
			// and receive acknowledge
			if err := game.Update(msg, func(ack string) {
				// do something ...
			}); err != nil {
				log.Println("Error: ", err)
			}
		}
	}
}
