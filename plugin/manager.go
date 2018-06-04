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
func Get(repo string, version string, local bool) (core.Game, error) {
	// Prepare plugin
	// Use local package when need to test game logical
	// If package being edited go get will error because of git
	// Otherwise it will re-download the package from github
	if !local {
		runcmd(fmt.Sprintf("go get -d -v -u %s", repo))
	}
	// Able to change version of package based on git
	// If no version the latest version of package
	if version != "" {
		runcmd(fmt.Sprintf("git -C $GOPATH/src/github.com/masaruz/engine-bomberman checkout %s", version))
	}
	runcmd(fmt.Sprintf("go build -buildmode=plugin -o $GOPATH/src/engine/%s $GOPATH/src/%s/main.go", file, repo))
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

func runcmd(cmd string) []byte {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		panic(err)
	}
	return out
}
