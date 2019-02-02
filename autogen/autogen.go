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
	cmake.FinalizeCompilerFlags()
	for name, buildinfo := range pkginfo.Targets.Libraries {
		cmake.BuildLibrary(name, buildinfo.Compile, buildinfo.Link)
	}
	for name, buildinfo := range pkginfo.Targets.Executables {
		cmake.BuildExecutable(name, buildinfo.Compile, buildinfo.Link)
	}
	for name, testInfo := range pkginfo.Tests {
		cmake.RunTest(name, testInfo.Command)
	}
}
