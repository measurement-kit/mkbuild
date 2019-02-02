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
			log.Warnf("unknown dependency: %s", depname)
			continue
		}
		handler(cmake)
	}
	cmake.SetRestrictiveCompilerFlags()
	cmake.WriteSectionComment("finalize compiler")
	cmake.WriteLine("add_definitions(${CMAKE_REQUIRED_DEFINITIONS})")
	cmake.WriteLine("include_directories(${CMAKE_REQUIRED_INCLUDES})")
	cmake.WriteLine("link_libraries(${CMAKE_REQUIRED_LIBRARIES})")
	cmake.WriteLine("enable_testing()")
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
