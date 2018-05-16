package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"

	"github.com/masaruz/engine-lib/core"
)

const (
	file = "main.so"
)

// Get plugin symbol
func Get(repo string) (core.Game, error) {
	// Prepare plugin
	execute(fmt.Sprintf("go get -d -v %s", repo))
	execute(fmt.Sprintf("go build -buildmode=plugin -o $GOPATH/src/engine/%s $GOPATH/src/%s/main.go", file, repo))
	////////////////////////////////////////////////////////
	//////// 1. Open the so file to load the symbols ///////
	////////////////////////////////////////////////////////
	path, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		return nil, err
	}
	dir := fmt.Sprintf("%s/%s", path, file)
	plug, err := plugin.Open(dir)
	if err != nil {
		return nil, err
	}
	////////////////////////////////////////////////////////////////////////
	//////// 2. Look up a symbol (an exported function or variable) ////////
	////////////////////////////////////////////////////////////////////////
	sym, err := plug.Lookup("Game")
	if err != nil {
		return nil, err
	}

	///////////////////////////////////////////////////////////////////
	//////// 3. Assert that loaded symbol is of a desired type ////////
	///////////////////////////////////////////////////////////////////
	game, ok := sym.(core.Game)
	if !ok {
		return nil, fmt.Errorf("unexpected type from module symbol")
	}
	// Remove .so after execute success
	os.Remove(dir)
	return game, nil
}

func execute(cmd string) {
	if _, err := exec.Command("bash", "-c", cmd).Output(); err != nil {
		panic(fmt.Sprintf("Failed to execute command: %s", cmd))
	}
}
