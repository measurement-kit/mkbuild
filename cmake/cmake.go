// Package cmake generates a CMakeLists.txt
package cmake

import (
	"sort"

	"github.com/apex/log"
	"github.com/measurement-kit/mkbuild/cmake/cmakefile"
	"github.com/measurement-kit/mkbuild/cmake/deps"
	"github.com/measurement-kit/mkbuild/pkginfo"
)

func sortedBuildInfo(m map[string]pkginfo.BuildInfo) []string {
	var res []string
	for k, _ := range m {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

func sortedLibraryBuildInfo(m map[string]pkginfo.LibraryBuildInfo) []string {
	var res []string
	for k, _ := range m {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

func sortedScriptBuildInfo(m map[string]pkginfo.ScriptBuildInfo) []string {
	var res []string
	for k, _ := range m {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

func sortedTestInfo(m map[string]pkginfo.TestInfo) []string {
	var res []string
	for k, _ := range m {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

// Generate generates a CMakeLists.txt file.
func Generate(pkginfo *pkginfo.PkgInfo) {
	cmake := cmakefile.Open(pkginfo.Name)
	defer cmake.Close()
	for _, depname := range pkginfo.Dependencies {
		handler, ok := deps.All[depname]
		if !ok {
			log.Fatalf("unknown dependency: %s", depname)
		}
		handler(cmake)
	}
	cmake.FinalizeCompilerFlags()
	for _, name := range sortedLibraryBuildInfo(pkginfo.Targets.Libraries) {
		buildinfo := pkginfo.Targets.Libraries[name]
		cmake.AddLibrary(
			name, buildinfo.Compile, buildinfo.Link, buildinfo.Install,
			buildinfo.Headers,
		)
	}
	for _, name := range sortedBuildInfo(pkginfo.Targets.Executables) {
		buildinfo := pkginfo.Targets.Executables[name]
		cmake.AddExecutable(
			name, buildinfo.Compile, buildinfo.Link, buildinfo.Install,
		)
	}
	for _, name := range sortedScriptBuildInfo(pkginfo.Targets.Scripts) {
		buildinfo := pkginfo.Targets.Scripts[name]
		cmake.AddScript(name, buildinfo.Install)
	}
	for _, name := range sortedTestInfo(pkginfo.Tests) {
		testinfo := pkginfo.Tests[name]
		cmake.AddTest(name, testinfo.Command)
	}
}
