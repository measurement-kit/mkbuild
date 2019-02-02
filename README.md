# MKBuild

MKBuild (1) generates a complex `CMakeLists.txt` from a simpler YAML based
definition of the build and (2) allows to easily run a specific build using
Docker to perform the build and testing in a specific container.

## Getting the software

```
go get -v github.com/bassosimone/mkbuild
```

## Autogenerating CMakeLists.txt

Create a YAML file named `MKBuild.yaml` in the toplevel directory of your
project and write inside it something similar to:

```YAML
name: mkcurl

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
  - tests
  integration_tests:
  - integration-tests
  external_ca:
  - mkcurl-client --ca-bundle-path ./.mkbuild/etc/ca-bundle.pem
    https://www.kernel.org
  http11_test:
  - mkcurl-client https://ooni.torproject.org
  using_timeout:
  - mkcurl-client --timeout 10 --follow-redirect https://www.facebook.com
  redirect_test:
  - mkcurl-client --follow-redirect http://google.com
  post:
  - mkcurl-client --post --data "{\"net-tests\":[]}" https://httpbin.org/post
  put:
  - mkcurl-client --put --data "{\"net-tests\":[]}" https://httpbin.org/put
  connect_to:
  - mkcurl-client --connect-to www.google.com https://www.youtube.com
```

Where `name` is the name of the project, `dependencies` is a list containing
the IDs of the dependencies you want to download and install, `targets` tells
us what artifacts you want to build, and `tests` what tests to execute.

See `autogen/rules/rules.go` for all the available IDs. Dependencies that
are libraries will be automatically downloaded for Windows, but must be
installed on Unix. If a dependency is not installed on Unix, the related
`cmake` check will fail when running `cmake` later on. The build flags will
be automatically adjusted to account for finding the dependencies headers
and for automatically linking all targets with dependencies.

The `libraries` key specifies what libraries to build and the `executables`
key what executables to build. Both contain targets names mapping to build
information for a target. The build information is composed of two keys,
`compile`, which indicates which sources to compile, and `link`, which
indicates which _extra_ libraries to link. If you're building, like in the
above example, a library named `foo`, you can refer to it later in the
`link` section of another target simply as `foo`.

The `tests` indicates what test to run. Each key inside `tests` is the name
of a test. Each key maps to a list of arguments to be passed to a test. It's
up to you whether to put each argument as a separate list item, or to put
all the arguments as part of the same list entry, like in figure. We do allow
for both styles, as the latter may be convenient with very long cmdlines.

One you've written you `MKBuild.yaml`, run

```
mkbuild autogen
```

This will generate a `CMakeLists.txt` file. From there on, just follow the
standard procedures to build with `cmake`. Note that dependencies will be
downloaded and configured by `cmake`, not by `mkbuild`, which just generates
a suitable `CMakeLists.txt` file to perform the task.

## Running a build in Docker

To run a build in docker, you should know about the type of builds that
are available. To this end, see `docker/docker.go`. The simples build
type is the `vanilla` build. Since this is a personal project, the docker
image that we'll use is the one used by Measurement Kit builds.

By running, e.g.

```
mkbuild docker vanilla
```

you will cause `mkbuild` to write a special bourne shell script in a
hidden directory and to launch `docker`, with the above mentioned docker
image, such that this script is run inside the container.

Such script will rebuild `mkbuild` inside the container and then use
it to perform the selected kind of build. This will basically boil down
to calling `mkbuild autogen` to generate a `CMakeLists.txt` and
then following the typical steps of a `cmake` build.

## Rationale

This software is meant to replace the `github.com/measurement-kit/cmake-utils`
and `github.com/measurement-kit/ci-scripts` subrepositories. Rather than
having to keep the submodules up to date, like we do, e.g., in `mkcurl`, one
`go get`s the latest `mkbuild` during a build to obtain the same result.

The main difference is that there is no need to keep in sync all the submodules
of the many small repositories I've created in `gitub.com/measurement-kit`. More
details in the following subsections.

Also, even in case I'm doing it wrong, it's still possible to cut
this tool of the build by commiting the `CMakeLists.txt`. Also,
in case we want to have ready-to-use tarballs for release (I doubt
it), we can generate a tarball with a `CMakeLists.txt` in it.

## Travis CI

The `.travis.yml` file should look like

```YAML
language: go

go:
- 1.11

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
  - go get -v github.com/bassosimone/mkbuild
  - $GOPATH/bin/mkbuild docker $BUILD_TYPE
```

It only minimally more complex than what was required by `ci-common`
and `cmake-modules`.

## AppVeyor

The `.appveyor.yml` file should look like:

```YAML
image: Visual Studio 2017

environment:
  GOPATH: c:/gopath
  GOVERSION: 1.11
  matrix:
    - CMAKE_GENERATOR: "Visual Studio 15 2017 Win64"
    - CMAKE_GENERATOR: "Visual Studio 15 2017"

build_script:
  - cmd: go get -v github.com/bassosimone/mkbuild
  - cmd: "%GOPATH%/bin/mkbuild.exe autogen"
  - cmd: cmake -G "%CMAKE_GENERATOR%"
  - cmd: cmake --build . -- /nologo /property:Configuration=Release
  - cmd: ctest --output-on-failure -C Release -a
```

It only minimally more complex than what was required by `ci-common`
and `cmake-modules`.

## Next steps

If testing proves that this is really more convenient, I will
most likely migrate it into the `measurement-kit` namespace.
