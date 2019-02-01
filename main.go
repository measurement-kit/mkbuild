package main

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/bassosimone/mkbuild/autogen"
	"github.com/bassosimone/mkbuild/docker"
)

func main() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
	if len(os.Args) == 2 && os.Args[1] == "autogen" {
		autogen.Run()
	} else if len(os.Args) == 3 && os.Args[1] == "docker" {
		docker.Run(os.Args[2])
	} else {
		fmt.Fprintf(os.Stderr, "Usage: mkbuild autogen\n")
		fmt.Fprintf(os.Stderr, "       mkbuild docker <build-type>\n")
		os.Exit(1)
	}
}
