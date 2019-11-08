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

// FunctionCheck adds a check for a specific function
type FunctionCheck struct {
	// Name is the function name
	Name string

	// Define is the define to add to the build if function exists
	Define string
}

// SymbolCheck adds a check for a specific symbol
type SymbolCheck struct {
	// Name is the symbol name
	Name string

	// Header is the header to include for checking for the symbol
	Header string

	// Define is the define to add to the build if symbol exists
	Define string
}

// PkgInfo contains information on a package
type PkgInfo struct {
	// Name is the name of the package
	Name string

	// FunctionChecks contains all the checks for functions
	FunctionChecks []FunctionCheck `yaml:"function_checks"`

	// SymbolChecks contains all the checks for symbols
	SymbolChecks []SymbolCheck `yaml:"symbol_checks"`

	// Docker is the docker container to use for running tests
	Docker string

	// DockerTcDisabled indicates whether we should disable using tc inside
	// the container to artificially increase the latency. This is usually
	// needed when you measure performance and is actually counterproductive
	// otherwise, because it slows down operations and may also cause some
	// more or less predicatable failures. The reason why this flag defaults
	// to not disabling `tc` is to provide backward compatibility.
	DockerTcDisabled bool `yaml:"docker_tc_disabled"`

	// Dependencies are the package dependencies
	Dependencies []string

	// Amalgamate maps names the name of an amalgamated file to the
	// sorted list of source files that should be amalgamated.
	Amalgamate map[string][]string

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
