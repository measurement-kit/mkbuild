// Package cmake generates a CMakeLists.txt
package cmake

import (
	"github.com/apex/log"
	"github.com/bassosimone/mkbuild/cmake/cmakefile"
	"github.com/bassosimone/mkbuild/cmake/deps"
	"github.com/bassosimone/mkbuild/pkginfo"
)

// Generate generates a CMakeLists.txt file.
func Generate(pkginfo *pkginfo.PkgInfo) {
	cmake := cmakefile.Open(pkginfo.Name)
	defer cmake.Close()
	for _, depname := range pkginfo.Dependencies {
		handler, ok := deps.Rules[depname]
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
