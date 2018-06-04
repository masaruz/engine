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
	if stage := os.Getenv("STAGE"); stage == "dev" {
		port = 3000
	} else {
		port, _ = freeport.GetFreePort()
	}

	game, err := plugin.Get(
		"github.com/masaruz/engine-bomberman",
		os.Getenv("TAG"),
		os.Getenv("LOCAL_PACKAGE") == "")

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
