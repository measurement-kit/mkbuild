package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/bassosimone/mkbuild/cmake"
	"github.com/bassosimone/mkbuild/pkginfo"
	"github.com/bassosimone/mkbuild/rules"
)

// subrAutogen implements the autogen behaviour.
func subrAutogen() {
	pkginfo := pkginfo.Read()
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
	rules.WriteSectionComment(cmake, "set restrictive compiler flags")
	cmake.SetRestrictiveCompilerFlags()
	rules.WriteSectionComment(cmake, "finalize compiler")
	cmake.WriteLine("add_definitions(${CMAKE_REQUIRED_DEFINITIONS})")
	cmake.WriteLine("include_directories(${CMAKE_REQUIRED_INCLUDES})")
	cmake.WriteLine("link_libraries(${CMAKE_REQUIRED_LIBRARIES})")
	cmake.WriteLine("enable_testing()")
	for name, sources := range pkginfo.Build.Executables {
		rules.WriteSectionComment(cmake, name)
		cmake.WriteLine(fmt.Sprintf("add_executable("))
		cmake.WriteLine(fmt.Sprintf("  %s", name))
		for _, source := range sources {
			cmake.WriteLine(fmt.Sprintf("  %s", source))
		}
		cmake.WriteLine(fmt.Sprintf(")"))
	}
	for name, arguments := range pkginfo.Tests {
		rules.WriteSectionComment(cmake, "test: " + name)
		cmake.WriteLine(fmt.Sprintf("add_test("))
		cmake.WriteLine(fmt.Sprintf("  NAME %s COMMAND", name))
		for _, arg := range arguments {
			cmake.WriteLine(fmt.Sprintf("  %s", arg))
		}
		cmake.WriteLine(fmt.Sprintf(")"))
	}
}

// runnerTemplate is the template runner.sh run in the container.
var runnerTemplate = `#!/bin/sh -e
BUILD_TYPE="{{.BUILD_TYPE}}"
CODECOV_TOKEN="{{.CODECOV_TOKEN}}"
TRAVIS_BRANCH="{{.TRAVIS_BRANCH}}"
set -x

# Build the latest mkbuild for the docker container
export GOPATH=/go
install -d $GOPATH
go get -v github.com/bassosimone/mkbuild
cd /mk
env|grep -v TOKEN|sort
$GOPATH/bin/mkbuild autogen

# Make sure we don't consume too much resources by bumping latency
tc qdisc add dev eth0 root netem delay 200ms 10ms

# Select the proper build flags depending on the build type
if [ "$BUILD_TYPE" = "asan" ]; then
  export CFLAGS="-fsanitize=address -O1 -fno-omit-frame-pointer"
  export CXXFLAGS="-fsanitize=address -O1 -fno-omit-frame-pointer"
  export LDFLAGS="-fsanitize=address -fno-omit-frame-pointer"
  export CMAKE_BUILD_TYPE="Debug"

elif [ "$BUILD_TYPE" = "clang" ]; then
  export CMAKE_BUILD_TYPE="Release"
  export CXXFLAGS="-stdlib=libc++"

elif [ "$BUILD_TYPE" = "coverage" ]; then
  export CFLAGS="-O0 -g -fprofile-arcs -ftest-coverage"
  export CMAKE_BUILD_TYPE="Debug"
  export CXXFLAGS="-O0 -g -fprofile-arcs -ftest-coverage"
  export LDFLAGS="-lgcov"

elif [ "$BUILD_TYPE" = "ubsan" ]; then
  export CFLAGS="-fsanitize=undefined -fno-sanitize-recover"
  export CXXFLAGS="-fsanitize=undefined -fno-sanitize-recover"
  export LDFLAGS="-fsanitize=undefined"
  export CMAKE_BUILD_TYPE="Debug"

elif [ "$BUILD_TYPE" = "vanilla" ]; then
  export CMAKE_BUILD_TYPE="Release"

else
  echo "$0: BUILD_TYPE not in: asan, clang, coverage, tsan, ubsan, vanilla" 1>&2
  exit 1
fi

# Configure, make, and make check equivalent
cmake -GNinja -DCMAKE_BUILD_TYPE=$CMAKE_BUILD_TYPE .
cmake --build . -- -v
ctest --output-on-failure -a -j8

# Measure and possibly report the test coverage
if [ "$BUILD_TYPE" = "coverage" ]; then
  lcov --directory . --capture -o lcov.info
  if [ "$CODECOV_TOKEN" != "" ]; then
    curl -fsSL -o codecov.sh https://codecov.io/bash
    bash codecov.sh -X gcov -f lcov.info
  fi
fi
`

// writeDockerRunner writes the docker runner script.
func writeDockerRunner(buildType string) {
	tmpl := template.Must(template.New("runner.sh").Parse(runnerTemplate))
	dirname := filepath.Join(".mkbuild", "script")
	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		log.WithError(err).Fatalf("cannot create dir: %s", dirname)
	}
	filename := filepath.Join(dirname, "runner.sh")
	filep, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		log.WithError(err).Fatalf("cannot open file: %s", filename)
	}
	defer filep.Close()
	err = tmpl.Execute(filep, map[string]string{
		"BUILD_TYPE":    buildType,
		"CODECOV_TOKEN": os.Getenv("CODECOV_TOKEN"),
		"TRAVIS_BRANCH": os.Getenv("TRAVIS_BRANCH"),
	})
	if err != nil {
		log.WithError(err).Fatalf("cannot write file: %s", filename)
	}
}

// subrDocker implements the docker behaviour.
func subrDocker(buildType string) {
	writeDockerRunner(buildType)
	cwd, err := os.Getwd()
	if err != nil {
		log.WithError(err).Fatal("os.Getwd failed")
	}
	command := exec.Command("docker", "run", "--cap-add=NET_ADMIN", "-v",
		fmt.Sprintf("%s:/mk", cwd), "-t", "bassosimone/mk-debian",
		"/mk/.mkbuild/script/runner.sh")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		log.WithError(err).Fatal("cannot run build inside docker")
	}
}

func main() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
	if len(os.Args) == 2 && os.Args[1] == "autogen" {
		subrAutogen()
	} else if len(os.Args) == 3 && os.Args[1] == "docker" {
		subrDocker(os.Args[2])
	} else {
		fmt.Fprintf(os.Stderr, "Usage: mkbuild autogen\n")
		fmt.Fprintf(os.Stderr, "       mkbuild docker <build-type>\n")
		os.Exit(1)
	}
}
