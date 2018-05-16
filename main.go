package main

import (
	"engine/network"
	"engine/plugin"
	"log"

	"github.com/phayes/freeport"
)

func main() {
	// get a free port
	port, _ := freeport.GetFreePort()

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
