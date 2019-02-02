// Package autogen allows to autogenerate a build system
package autogen

import (
	"fmt"

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
		cmake.WriteSectionComment(name)
		cmake.WriteLine(fmt.Sprintf("add_executable("))
		cmake.WriteLine(fmt.Sprintf("  %s", name))
		for _, source := range sources {
			cmake.WriteLine(fmt.Sprintf("  %s", source))
		}
		cmake.WriteLine(fmt.Sprintf(")"))
	}
	for name, arguments := range pkginfo.Tests {
		cmake.WriteSectionComment("test: "+name)
		cmake.WriteLine(fmt.Sprintf("add_test("))
		cmake.WriteLine(fmt.Sprintf("  NAME %s COMMAND", name))
		for _, arg := range arguments {
			cmake.WriteLine(fmt.Sprintf("  %s", arg))
		}
		cmake.WriteLine(fmt.Sprintf(")"))
	}
}
