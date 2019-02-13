package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/measurement-kit/mkbuild/cmake"
	"github.com/measurement-kit/mkbuild/docker"
	"github.com/measurement-kit/mkbuild/pkginfo"
)

func main() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
	pkginfo := pkginfo.Read()
	docker.Generate(pkginfo)
	cmake.Generate(pkginfo)
}
