package main

import (
	"os"
        _ "embed"

	"github.com/aptly-dev/aptly/aptly"
	"github.com/aptly-dev/aptly/cmd"
)

// Version variable, filled in at link time
//go:generate sh -c "make -s version | tr -d '\n' > VERSION"
//go:embed VERSION
var Version string

func main() {
	if Version == "" {
		Version = "unknown"
	}

	aptly.Version = Version

	os.Exit(cmd.Run(cmd.RootCommand(), os.Args[1:], true))
}
