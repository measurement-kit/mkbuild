# MKBuild

MKBuild is a small utility to simplify managing Measurement
Kit builds. This tool performs two main tasks:

1. generates or updates a `CMakeLists.txt` that downloads
   required dependencies, configures strict compiler flags,
   builds libraries and executables, and run tests;

2. generates or updates the `docker.sh` that runs a CMake
   based build inside a specific Docker container, with
   specific compiler flags (e.g. for `asan`).

MKBuild is driven by the configuration contained in the
`MKBuild.yml` YAML file. Read on for more info.

## Getting MKBuild

```
go get -v github.com/measurement-kit/mkbuild
```

## Converting a repository to use MKBuild

Create `MKBuild.yaml` in the toplevel directory of your project. This
file should look like this:

```YAML
name: mkcurl

docker: bassosimone/mk-debian

dependencies:
- curl.haxx.se/ca
- github.com/adishavit/argh
- github.com/catchorg/catch2
- github.com/curl/curl
- github.com/measurement-kit/mkmock

targets:
  libraries:
    mkcurl:
      compile: [mkcurl.cpp]
  executables:
    mkcurl-client:
      compile: [mkcurl-client.cpp]
      link: [mkcurl]
    tests:
      compile: [tests.cpp]
    integration-tests:
      compile: [integration-tests.cpp]
      link: [mkcurl]

tests:
  mocked_tests:
    command: tests
  integration_tests:
    command: integration-tests
  redirect_test:
    command: mkcurl-client --follow-redirect http://google.com
```

Where `name` is the name of the project, `docker` is the name of the
docker container to use, `dependencies` lists the IDs of the dependencies
you want to download and install, `targets` tells us what artifacts you
want to build, and `tests` what tests to execute.

See `cmake/deps/deps.go` for all the available deps IDs. Dependencies
that compile to static/shared libraries (e.g. `libcurl`) will be downloaded
automatically on Windows, and must be already installed on Unix systems. If a
dependency is not already installed on Unix, the related `cmake` check will
fail when running `cmake` later on. The build flags will be automatically
adjusted to account for a dependency (e.g. `CXXFLAGS` and `LDFLAGS` will be
linked to use cURL's headers and libraries).

The `libraries` key specifies what libraries to build and the
`executables` key what executables to build. Both contain targets names
mapping to build information for a target. The build information is
composed of two keys, `compile`, which indicates which sources to compile,
and `link`, which indicates which libraries to link. You do not need to
list here the dependencies, but you can list here libraries built as part
of the local build. In the above example, the `integration-tests` binary
will link with the (static) library called `mkcurl`, in addition to linking
to all the libraries implied by the declared dependencies.

The `tests` key indicates what test to run. Each key inside `tests` is the name
of a test. The `command` key indicates what command to execute.

## (Re)Generating CMakeLists.txt and docker.sh

One you've written (or updated) `MKBuild.yaml`, just run

```
mkbuild
```

This will generate (or update) `CMakeLists.txt` and `docker.sh`.

You should commit these files to the repository.

## Build instructions

Since `mkbuild` generates a `CMakeLists.txt` and we suggest to commit
it to your repository, the build instructions are the standard build
instructions of any CMake based software project.

## Running a build using Docker

Provided that you have Docker installed, running a docker based
build is as simple as running:

```
./docker.sh <build-type>
```

Run `docker.sh` without arguments to see the available build types. The
names of the build types should be self explanatory.

## Rationale

This software is meant to replace the `github.com/measurement-kit/cmake-utils`
and `github.com/measurement-kit/ci-scripts` subrepositories. Rather than
having to keep the submodules up to date, we automatically generate files
and scripts implementing the same functionality.

Because this tool generates standalone `CMakeLists.txt` and `docker.sh`, it
means that it can easily be replaced with better tools, or no tools. Yet, the
burden of keeping in sync the subrepos is gone and it is replaced with the
much lower burden of running `mkbuild` from time to time to sync.

An earlier design of this tool was such that `CMakeLists.txt` and `docker.sh`
were not committed to the repository. Yet, this is probably not advisable since
it may lead to non reproducible continuous integration builds, because the
newly generated `CMakeLists.txt` and/or `docker.sh` may differ. In any case,
should we decided that _not committing_ these files into the repository is
instead better, we just need to update the build instructions to mention to
compile and run `mkbuild` as the first step.

## Travis CI

The `.travis.yml` file should look like:

```YAML
language: c++
services:
  - docker
sudo: required
matrix:
  include:
    - env: BUILD_TYPE="asan"
    - env: BUILD_TYPE="clang"
    - env: BUILD_TYPE="coverage"
    - env: BUILD_TYPE="ubsan"
    - env: BUILD_TYPE="vanilla"
script:
  - ./docker.sh $BUILD_TYPE
```

This is equal to what we have now, _except_ that the name of the script
differs from the `github.com/measurement-kit/ci-common` one.

## AppVeyor

The `.appveyor.yml` is quite like the one that we use now:

```YAML
image: Visual Studio 2017
environment:
  matrix:
    - CMAKE_GENERATOR: "Visual Studio 15 2017 Win64"
    - CMAKE_GENERATOR: "Visual Studio 15 2017"
build_script:
  - cmd: cmake -G "%CMAKE_GENERATOR%"
  - cmd: cmake --build . -- /nologo /property:Configuration=Release
  - cmd: ctest --output-on-failure -C Release -a
```

The main difference is that we don't need to update subrepos anymore.
