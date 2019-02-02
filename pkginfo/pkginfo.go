// Package pkginfo contains information on a package
package pkginfo

import (
	"io/ioutil"

	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

type buildInfo struct {
	// Compile lists all the sources to compile
	Compile []string

	// Link lists all the libraries to link
	Link []string
}

type targetInfo struct {
	// Libraries lists all the libraries to build
	Libraries map[string]buildInfo

	// Executables lists all the executabls to build
	Executables map[string]buildInfo
}

// PkgInfo contains information on a package
type PkgInfo struct {
	// Name is the name of the package
	Name string

	// Dependencies are the package dependencies
	Dependencies []string

	// Targets contains information on what we need to build
	Targets targetInfo

	// Tests contains information on the tests to run
	Tests map[string][]string
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
