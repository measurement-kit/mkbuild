# MKBuild

MKBuild (1) generates a `CMakeLists.txt` and (2) scripts for testing on
a Docker container from a simple YAML based definition.

## Getting the software

```
go get -v github.com/bassosimone/mkbuild
```

## (Re)generating build and test scripts.

Create (or update) a YAML file named `MKBuild.yaml` in the toplevel
directory of your project. This file should look like this:

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
that are libraries will be automatically downloaded for Windows, but
must be installed on Unix. If a dependency is not installed on Unix,
the related `cmake` check will fail when running `cmake` later on. The
build flags will be automatically adjusted to account for compiling and
linking with the specified dependencies.

The `libraries` key specifies what libraries to build and the
`executables` key what executables to build. Both contain targets names
mapping to build information for a target. The build information is
composed of two keys, `compile`, which indicates which sources to compile,
and `link`, which indicates which _extra_ libraries to link. (Remember
that dependencies will automatically be accounted for, so you don't
need to say you want to link with them explicitly.) If you're building,
like in the above example, a library named `foo`, you can refer to it
later in the `link` section of another target simply as `foo`.

The `tests` indicates what test to run. Each key inside `tests` is the name
of a test. The `command` key indicates what command to execute.

One you've written (or updated) your `MKBuild.yaml`, run

```
mkbuild
```

This will generate (or update) thee files:

1. `CMakeLists.txt`

2. `.ci/docker/trampoline.sh`

3. `.ci/docker/run.sh`

You should commit these files to the repository.

The `CMakeLists.txt` file will contain the rules to download and use
the dependencies, build the required artifacts, and run the tests. Every
run of `mkbuild` updates this file with the latest known version of the
dependencies. That is, we don't support version pinning, as we aim to live
as close to the latest version of the dependencies as possible.

The `.ci/docker/trampoline.sh` contains the code to run the `run.sh`
script inside a specific Docker container. The `run.sh` script contains
code to run several kind of Linux based, CMake based builds, including
for example `asan`, `tsan`, and coverage builds.

## Build instructions

Since `mkbuild` generates a `CMakeLists.txt` and we suggest to commit
it to your repository, the build instructions are the standard build
instructions of any CMake based software project.

## Running a build using Docker

Provided that you have Docker installed, running a docker based
build is as simple as running:

```
./.ci/docker/trampoline.sh <build-type>
```

Run `trampoline.sh` without arguments to see the available build types.

## Rationale

This software is meant to replace the `github.com/measurement-kit/cmake-utils`
and `github.com/measurement-kit/ci-scripts` subrepositories. Rather than
having to keep the submodules up to date, we automatically generate files
and scripts implementing the same functionality.

Because this tool generates standalone `CMakeLists.txt` and shell scripts, it
means that it can easily be replaced with better tools, or no tools, in the
future, without any annoyance. Yet, the burden of keeping in sync the subrepos
is gone and it is replaced with the much lower burden of running `mkbuild`
from time to time to stay in sync.

An earlier design of this tool was such that `CMakeLists.txt` and the scripts
were not committed to the repository. Yet, this is probably not advisable as
it may lead to non reproducible continuous integration builds, because the
newly generated CMakeLists.txt or scripts may differ. In any case, should we
decided that _not committing_ `CMakeLists.txt` and the scripts into the
repository is instead better, we just need to update the build instructions
to mention to compile and run `mkbuild` as the first step.

## Travis CI

The `.travis.yml` file should look like

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
  - ./.ci/docker/trampoline.sh $BUILD_TYPE
```

This is equal to what we have now, _except_ that the name of the script
differs from the `github.com/measurement-kit/ci-common` one.

## AppVeyor

The `.appveyor.yml` is like the one that we use now, that is:

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

## Next steps

If testing proves that this repository is really more convenient, I will
most likely migrate it into the `measurement-kit` namespace.
