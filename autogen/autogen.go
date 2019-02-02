// Package autogen allows to autogenerate a build system
package autogen

import (
	"github.com/apex/log"
	"github.com/bassosimone/mkbuild/autogen/cmake"
	"github.com/bassosimone/mkbuild/autogen/rules"
	"github.com/bassosimone/mkbuild/pkginfo"
)

// Run implements the autogen subcommand.
func Run(pkginfo *pkginfo.PkgInfo) {
	cmake := cmake.Open(pkginfo.Name)
	defer cmake.Close()
	for _, depname := range pkginfo.Dependencies {
		handler, ok := rules.Rules[depname]
		if !ok {
			log.Fatalf("unknown dependency: %s", depname)
		}
		handler(cmake)
	}
	cmake.SetRestrictiveCompilerFlags()
	cmake.PrepareForCompilingTargets()
	for name, sources := range pkginfo.Build.Executables {
		cmake.AddExecutable(name, sources)
	}
	for name, arguments := range pkginfo.Tests {
		cmake.AddTest(name, arguments)
	}
}
