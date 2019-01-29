package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"gopkg.in/yaml.v2"
)

// download will download |URL| in |filename|.
func download(filename, URL string) {
	dirname, _ := filepath.Split(filename)
	if dirname != "" {
		err := os.MkdirAll(dirname, 0755)
		if err != nil {
			log.WithError(err).Fatalf("MkdirAll failed for: %s", dirname)
		}
	}
	filep, err := os.Create(filename)
	if err != nil {
		log.WithError(err).Fatalf("os.Create failed for: %s", filename)
	}
	defer filep.Close()
	response, err := http.Get(URL)
	if err != nil {
		log.WithError(err).Fatalf("http.Get failed for: %s", URL)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatalf("HTTP server response: %s; for URL: %s", response.Status, URL)
	}
	_, err = io.Copy(filep, response.Body)
	if err != nil {
		log.WithError(err).Fatalf("io.Copy failed for: %s", filename)
	}
}

// verify will verify that |filename| has SHA256 equal to |SHA256|.
func verify(filename, SHA256 string) {
	filep, err := os.Open(filename)
	if err != nil {
		log.WithError(err).Fatalf("log.Open failed for: %s", filename)
	}
	defer filep.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, filep); err != nil {
		log.WithError(err).Fatalf("io.Copy failed for: %s", filename)
	}
	result := hex.EncodeToString(hash.Sum(nil))
	if result != SHA256 {
		log.Fatalf("hash mismatch for: %s; sha256: %s", filename, SHA256)
	}
}

// downloadAndVerify downloads |URL| in |filename| and verify that
// |filename| has SHA256 equal to |SHA256|.
func downloadAndVerify(filename, SHA256, URL string) {
	download(filename, URL)
	verify(filename, SHA256)
}

// moduleInfo contains info on the module
type moduleInfo struct {
	// Name is the name of the module read from MKBuild.yaml
	Name string `yaml:"name"`

	// Dependencies are the module dependencies read from MKBuild.yaml
	Dependencies []string `yaml:"dependencies"`

	// Executables lists all the executabls to build read from MKBuild.yaml
	Executables map[string][]string

	// Tests contains information on the tests to run read from MKBuild.yaml
	Tests map[string][]string

	// IncludeDirs are the include directories computed by the code
	// that installs all the dependencies
	IncludeDirs []string

	// IncludeDirsStr is the list of include directories formatted for
	// cmake as computed by the code that writes CMakeLists.txt
	IncludeDirsStr string

	// LinkLibs are the link libraries computed by the code
	// that installs all the dependencies
	LinkLibs []string

	// LinkLibsStr is the list of link libraries formatted for
	// cmake as computed by the code that writes CMakeLists.txt
	LinkLibsStr string
}

// gModuleInfo is the global moduleInfo
var gModuleInfo moduleInfo

// installCurlHaxxSeCa installs CURL's CA bundle
func installCurlHaxxSeCa(dep string) {
	log.Infof("install: %s", dep)
	downloadAndVerify(
		filepath.Join(".mkbuild", "dep", "curl.haxx.se", "ca", "ca-bundle.pem"),
		"4d89992b90f3e177ab1d895c00e8cded6c9009bec9d56981ff4f0a59e9cc56d6",
		"https://curl.haxx.se/ca/cacert-2018-12-05.pem",
	)
}

// installGithubcomAdishavitArgh installs github.com/adishavit/argh
func installGithubcomAdishavitArgh(dep string) {
	log.Infof("install: %s", dep)
	downloadAndVerify(
		filepath.Join(".mkbuild", "dep", "github.com", "adishavit", "argh", "argh.h"),
		"ddb7dfc18dcf90149735b76fb2cff101067453a1df1943a6911233cb7085980c",
		"https://raw.githubusercontent.com/adishavit/argh/v1.3.0/argh.h",
	)
	gModuleInfo.IncludeDirs = append(gModuleInfo.IncludeDirs,
		filepath.Join(".mkbuild", "dep", "github.com", "adishavit", "argh"))
}

// installGithubcomCatchorgCatch2 installs github.com/catchorg/Catch2
func installGithubcomCatchorgCatch2(dep string) {
	log.Infof("install: %s", dep)
	downloadAndVerify(
		filepath.Join(".mkbuild", "dep", "github.com", "catchorg", "Catch2", "catch.hpp"),
		"5eb8532fd5ec0d28433eba8a749102fd1f98078c5ebf35ad607fb2455a000004",
		"https://github.com/catchorg/Catch2/releases/download/v2.3.0/catch.hpp",
	)
	gModuleInfo.IncludeDirs = append(gModuleInfo.IncludeDirs,
		filepath.Join(".mkbuild", "dep", "github.com", "catchorg", "Catch2"))
}

// installGithubcomCurlCurl installs github.com/curl/curl
func installGithubcomCurlCurl(dep string) {
	log.Infof("install: %s", dep)
	gModuleInfo.LinkLibs = append(gModuleInfo.LinkLibs, "-lcurl")
}

// installGithubcomMeasurementkitMkmock installs
// github.com/measurement-kit/mkmock
func installGithubcomMeasurementkitMkmock(dep string) {
	log.Infof("install: %s", dep)
	downloadAndVerify(
		filepath.Join(".mkbuild", "dep", "github.com", "measurement-kit", "mkmock", "mkmock.hpp"),
		"f07bc063a2e64484482f986501003e45ead653ea3f53fadbdb45c17a51d916d2",
		"https://raw.githubusercontent.com/measurement-kit/mkmock/v0.2.0/mkmock.hpp",
	)
	gModuleInfo.IncludeDirs = append(gModuleInfo.IncludeDirs,
		filepath.Join(".mkbuild", "dep", "github.com", "measurement-kit", "mkmock"))
}

// cmakeTemplate is the template for CMakeLists.txt
var cmakeTemplate = `# Autogenerated by mkbuild
cmake_minimum_required(VERSION 3.1.0)
project({{.Name}})

set(CMAKE_POSITION_INDEPENDENT_CODE ON)
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_C_EXTENSIONS OFF)

set(THREADS_PREFER_PTHREAD_FLAG ON)
find_package(Threads REQUIRED)

if(("${CMAKE_CXX_COMPILER_ID}" STREQUAL "GNU") OR
   ("${CMAKE_CXX_COMPILER_ID}" MATCHES "Clang"))
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Werror")
  # https://www.owasp.org/index.php/C-Based_Toolchain_Hardening_Cheat_Sheet
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wall")
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wextra")
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wconversion")
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wcast-align")
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wformat=2")
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wformat-security")
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -fno-common")
  # Some options are only supported by GCC when we're compiling C code:
  if ("${CMAKE_CXX_COMPILER_ID}" MATCHES "Clang")
    set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wmissing-prototypes")
    set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wstrict-prototypes")
  else()
    set(MK_C_FLAGS "${MK_C_FLAGS} -Wmissing-prototypes")
    set(MK_C_FLAGS "${MK_C_FLAGS} -Wstrict-prototypes")
  endif()
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wmissing-declarations")
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wstrict-overflow")
  if("${CMAKE_CXX_COMPILER_ID}" STREQUAL "GNU")
    set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wtrampolines")
  endif()
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Woverloaded-virtual")
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wreorder")
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wsign-promo")
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -Wnon-virtual-dtor")
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} -fstack-protector-all")
  if(NOT "${APPLE}")
    set(MK_LD_FLAGS "${MK_LD_FLAGS} -Wl,-z,noexecstack")
    set(MK_LD_FLAGS "${MK_LD_FLAGS} -Wl,-z,now")
    set(MK_LD_FLAGS "${MK_LD_FLAGS} -Wl,-z,relro")
    set(MK_LD_FLAGS "${MK_LD_FLAGS} -Wl,-z,nodlopen")
    set(MK_LD_FLAGS "${MK_LD_FLAGS} -Wl,-z,nodump")
  endif()
  add_definitions(-D_FORTIFY_SOURCES=2)
elseif("${CMAKE_CXX_COMPILER_ID}" STREQUAL "MSVC")
  # TODO(bassosimone): add support for /Wall and /analyze
  set(MK_COMMON_FLAGS "${MK_COMMON_FLAGS} /WX /W4")
  set(MK_LD_FLAGS "${MK_LD_FLAGS} /WX")
else()
  message(FATAL_ERROR "Compiler not supported: ${CMAKE_CXX_COMPILER_ID}")
endif()
set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} ${MK_COMMON_FLAGS} ${MK_C_FLAGS}")
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} ${MK_COMMON_FLAGS} ${MK_CXX_FLAGS}")
set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} ${MK_LD_FLAGS}")
set(CMAKE_SHARED_LINKER_FLAGS "${CMAKE_SHARED_LINKER_FLAGS} ${MK_LD_FLAGS}")
if("${WIN32}")
  add_definitions(-D_WIN32_WINNT=0x0600) # for NI_NUMERICSERV and WSAPoll
endif()

include_directories({{.IncludeDirsStr}})

set(MK_LINK_LIBS {{.LinkLibsStr}})
if("${WIN32}" OR "${MINGW}")
  list(APPEND MK_LINK_LIBS "ws2_32")
  if ("${MINGW}")
      list(APPEND MK_LINK_LIBS -static-libgcc -static-libstdc++)
  endif()
endif()
list(APPEND MK_LINK_LIBS Threads::Threads)
link_libraries("${MK_LINK_LIBS}")

if("${WIN32}")
  compile_options(PRIVATE /EHs) # exceptions in extern "C"
endif()

enable_testing()
{{range $exeName, $sources := .Executables}}
add_executable({{$exeName}}{{range $idx, $src := $sources}} {{$src}}{{end}}){{end}}
{{range $testName, $cmdLine := .Tests}}
add_test(NAME {{$testName}} COMMAND {{range $idx, $arg := $cmdLine}}{{$arg}}{{end}}){{end}}
`

// writeCMakeListsTxt writes CMakeLists.txt in the current directory.
func writeCMakeListsTxt() {
	gModuleInfo.IncludeDirsStr = strings.Join(gModuleInfo.IncludeDirs, ";")
	gModuleInfo.LinkLibsStr = strings.Join(gModuleInfo.LinkLibs, ";")
	tmpl := template.Must(template.New("CMakeLists.txt").Parse(cmakeTemplate))
	filename := "CMakeLists.txt"
	filep, err := os.Create(filename)
	if err != nil {
		log.WithError(err).Fatalf("os.Open failed for: %s", filename)
	}
	defer filep.Close()
	err = tmpl.Execute(filep, gModuleInfo)
	if err != nil {
		log.WithError(err).Fatalf("tmpl.Execute failed for: %s", filename)
	}
	log.Infof("Written %s", filename)
}

// initializeModuleInfo reads module info from MKBuild.toml
func initializeModuleInfo() {
	filename := "MKBuild.yaml"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithError(err).Fatalf("cannot read %s", filename)
	}
	err = yaml.Unmarshal(data, &gModuleInfo)
	if err != nil {
		log.WithError(err).Fatalf("cannot unmarshal %s", filename)
	}
}

// satisfyDeps satisfies the dependencies
func satisfyDeps() {
	for _, dep := range gModuleInfo.Dependencies {
		if dep == "curl.haxx.se/ca" {
			installCurlHaxxSeCa(dep)
		} else if dep == "github.com/adishavit/argh" {
			installGithubcomAdishavitArgh(dep)
		} else if dep == "github.com/catchorg/catch2" {
			installGithubcomCatchorgCatch2(dep)
		} else if dep == "github.com/curl/curl" {
			installGithubcomCurlCurl(dep)
		} else if dep == "github.com/measurement-kit/mkmock" {
			installGithubcomMeasurementkitMkmock(dep)
		} else {
			log.Fatalf("unknown dependency: %s", dep)
		}
	}
}

// subrAutogen implements the autogen behaviour.
func subrAutogen() {
	initializeModuleInfo()
	satisfyDeps()
	writeCMakeListsTxt()
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
