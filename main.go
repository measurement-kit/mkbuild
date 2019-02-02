package main

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/bassosimone/mkbuild/autogen"
	"github.com/bassosimone/mkbuild/docker"
	"github.com/bassosimone/mkbuild/pkginfo"
)

func main() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
	pkginfo := pkginfo.Read()
	if len(os.Args) == 2 && os.Args[1] == "autogen" {
		autogen.Run(pkginfo)
	} else if len(os.Args) == 3 && os.Args[1] == "docker" {
		docker.Run(pkginfo, os.Args[2])
	} else {
		fmt.Fprintf(os.Stderr, "Usage: mkbuild autogen\n")
		fmt.Fprintf(os.Stderr, "       mkbuild docker <build-type>\n")
		os.Exit(1)
	}
}
