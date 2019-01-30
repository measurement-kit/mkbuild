// Package pkginfo contains information on a package
package pkginfo

import (
	"io/ioutil"

	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

type buildInfo struct {
	// Executables lists all the executabls to build read from MKBuild.yaml
	Executables map[string][]string
}

// PkgInfo contains information on a package
type PkgInfo struct {
	// Name is the name of the package
	Name string `yaml:"name"`

	// Dependencies are the module dependencies
	Dependencies []string `yaml:"dependencies"`

	// Build contains information on what we need to build
	Build buildInfo `yaml:"build"`

	// Tests contains information on the tests to run
	Tests map[string][]string `yaml:"tests"`
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
