package main

import (
	"engine/network"
	"engine/plugin"
	"log"
	"os"

	"github.com/phayes/freeport"
)

func main() {
	// get a free port
	var port int
	if stage := os.Getenv("STAGE"); stage == "prd" {
		port, _ = freeport.GetFreePort()
	} else {
		port = 3000
	}

	game, err := plugin.Get("github.com/masaruz/engine-bomberman")

	if err != nil {
		panic(err)
	}

	if err := game.Init(); err != nil {
		panic(err)
	}

	if err := game.Start(); err != nil {
		panic(err)
	}

	log.Fatal(network.ListenAndServe(port, game))
}
