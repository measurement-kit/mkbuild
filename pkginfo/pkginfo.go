// Package pkginfo contains information on a package
package pkginfo

import (
	"io/ioutil"

	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

// BuildInfo contains info on building a target
type BuildInfo struct {
	// Compile lists all the sources to compile
	Compile []string

	// Link lists all the libraries to link
	Link []string

	// Install indicates whether to install the target
	Install bool
}

// LibraryBuildInfo contains info on building a library
type LibraryBuildInfo struct {
	// Compile lists all the sources to compile
	Compile []string

	// Link lists all the libraries to link
	Link []string

	// Install indicates whether to install the library
	Install bool

	// Headers contains all the public headers
	Headers []string
}

// ScriptBuildInfo contains info on building a script
type ScriptBuildInfo struct {
	// Install indicates whether to install the script
	Install bool
}

// TargetsInfo contains info on all targets
type TargetsInfo struct {
	// Libraries lists all the libraries to build
	Libraries map[string]LibraryBuildInfo

	// Executables lists all the executables to build
	Executables map[string]BuildInfo

	// Scripts lists all the scripts to build
	Scripts map[string]ScriptBuildInfo
}

// TestInfo contains info on a test
type TestInfo struct {
	// Command is the command to execute
	Command string
}

// PkgInfo contains information on a package
type PkgInfo struct {
	// Name is the name of the package
	Name string

	// Docker is the docker container to use for running tests
	Docker string

	// Dependencies are the package dependencies
	Dependencies []string

	// Targets contains information on what we need to build
	Targets TargetsInfo

	// Tests contains information on the tests to run
	Tests map[string]TestInfo
}

// Read reads package info from "MKBuild.yaml"
func Read() *PkgInfo {
	filename := "MKBuild.yaml"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithError(err).Fatalf("cannot read %s", filename)
	}
	pkginfo := &PkgInfo{}
	err = yaml.Unmarshal(data, pkginfo)
	if err != nil {
		log.WithError(err).Fatalf("cannot unmarshal %s", filename)
	}
	return pkginfo
}
