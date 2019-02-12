package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/bassosimone/mkbuild/cmake"
	"github.com/bassosimone/mkbuild/docker"
	"github.com/bassosimone/mkbuild/pkginfo"
)

func main() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
	pkginfo := pkginfo.Read()
	docker.Generate(pkginfo)
	cmake.Generate(pkginfo)
}
