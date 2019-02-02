// Package cmake implements the CMake driver
package cmake

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/bassosimone/mkbuild/autogen/cmake/restrictiveflags"
)

// CMake is the CMake driver
type CMake struct {
	// output contains the CMakeLists.txt lines
	output strings.Builder

	// indent is the indent string to prefix to each line
	indent string
}

// WithIndent runs |func| with the specified |indent|.
func (cmake *CMake) WithIndent(indent string, fn func()) {
	oldIndent := cmake.indent
	cmake.indent += indent
	fn()
	cmake.indent = oldIndent
}

// WriteSectionComment writes a comment for |name| in |cmake|.
func (cmake *CMake) WriteSectionComment(name string) {
	cmake.WriteEmptyLine()
	cmake.WriteLine(fmt.Sprintf("#"))
	cmake.WriteLine(fmt.Sprintf("# %s", name))
	cmake.WriteLine(fmt.Sprintf("#"))
	cmake.WriteEmptyLine()
}

// WriteEmptyLine writes an empty line to output.
func (cmake *CMake) WriteEmptyLine() {
	cmake.WriteLine("")
}

// WriteLine writes a line to the CMakeLists.txt file.
func (cmake *CMake) WriteLine(s string) {
	if s != "" {
		_, err := cmake.output.WriteString(cmake.indent)
		if err != nil {
			log.WithError(err).Fatal("cannot write indent")
		}
		_, err = cmake.output.WriteString(s)
		if err != nil {
			log.WithError(err).Fatal("cannot write string")
		}
	}
	_, err := cmake.output.WriteString("\n")
	if err != nil {
		log.WithError(err).Fatal("cannot write newline")
	}
}

// Open opens a CMake project named |name|.
func Open(name string) *CMake {
	cmake := &CMake{}
	cmake.WriteLine("# Autogenerated file; DO NOT EDIT!")
	cmake.WriteLine(fmt.Sprintf("cmake_minimum_required(VERSION 3.12.0)"))
	cmake.WriteLine(fmt.Sprintf("project(\"%s\")", name))
	cmake.WriteEmptyLine()
	cmake.WriteLine("include(CheckIncludeFileCXX)")
	cmake.WriteLine("include(CheckLibraryExists)")
	cmake.WriteLine("include(CheckCXXCompilerFlag)")
	cmake.WriteLine("set(THREADS_PREFER_PTHREAD_FLAG ON)")
	cmake.WriteLine("find_package(Threads REQUIRED)")
	cmake.WriteLine("set(CMAKE_POSITION_INDEPENDENT_CODE ON)")
	cmake.WriteLine("set(CMAKE_CXX_STANDARD 11)")
	cmake.WriteLine("set(CMAKE_CXX_STANDARD_REQUIRED ON)")
	cmake.WriteLine("set(CMAKE_CXX_EXTENSIONS OFF)")
	cmake.WriteLine("set(CMAKE_C_STANDARD 11)")
	cmake.WriteLine("set(CMAKE_C_STANDARD_REQUIRED ON)")
	cmake.WriteLine("set(CMAKE_C_EXTENSIONS OFF)")
	cmake.WriteLine("list(APPEND CMAKE_REQUIRED_LIBRARIES Threads::Threads)")
	cmake.WriteLine("if(\"${WIN32}\")")
	cmake.WriteLine("  list(APPEND CMAKE_REQUIRED_LIBRARIES ws2_32 crypt32)")
	cmake.WriteLine("  if(\"${MINGW}\")")
	cmake.WriteLine("    list(APPEND CMAKE_REQUIRED_LIBRARIES -static-libgcc -static-libstdc++)")
	cmake.WriteLine("  endif()")
	cmake.WriteLine("endif()")
	cmake.WriteLine("enable_testing()")
	return cmake
}

// Download downloads |URL| to |filename| and checks the |SHA256|.
func (cmake *CMake) Download(filename, SHA256, URL string) {
	cmake.WriteLine(fmt.Sprintf("message(STATUS \"Download: %s\")", URL))
	cmake.WriteLine(fmt.Sprintf("file(DOWNLOAD %s", URL))
	cmake.WriteLine(fmt.Sprintf("  \"%s\"", filename))
	cmake.WriteLine(fmt.Sprintf("  EXPECTED_HASH SHA256=%s", SHA256))
	cmake.WriteLine(fmt.Sprintf("  TLS_VERIFY ON)"))
	cmake.WriteEmptyLine()
}

// checkCommandError writes the code to check for errors after a
// command has been executed.
func (cmake *CMake) checkCommandError() {
	cmake.WriteLine(fmt.Sprintf("if(\"${FAILURE}\")"))
	cmake.WriteLine(fmt.Sprintf("  message(FATAL_ERROR \"${FAILURE}\")"))
	cmake.WriteLine(fmt.Sprintf("endif()"))
}

// MkdirAll creates |destdirs|.
func (cmake *CMake) MkdirAll(destdirs string) {
	cmake.WriteLine(fmt.Sprintf("message(STATUS \"MkdirAll: %s\")", destdirs))
	cmake.WriteLine(fmt.Sprintf("execute_process(COMMAND"))
	cmake.WriteLine(fmt.Sprintf(
		"  ${CMAKE_COMMAND} -E make_directory \"%s\"", destdirs,
	))
	cmake.WriteLine(fmt.Sprintf("  RESULT_VARIABLE FAILURE)"))
	cmake.checkCommandError()
	cmake.WriteEmptyLine()
}

// Unzip extracts |filename| in |destdir|.
func (cmake *CMake) Unzip(filename, destdir string) {
	cmake.WriteLine(fmt.Sprintf("message(STATUS \"Extract: %s\")", filename))
	cmake.WriteLine(fmt.Sprintf("execute_process(COMMAND"))
	cmake.WriteLine(fmt.Sprintf(
		"  ${CMAKE_COMMAND} -E tar xf \"%s\"", filename,
	))
	cmake.WriteLine(fmt.Sprintf("  WORKING_DIRECTORY \"%s\"", destdir))
	cmake.WriteLine(fmt.Sprintf("  RESULT_VARIABLE FAILURE)"))
	cmake.checkCommandError()
	cmake.WriteEmptyLine()
}

// Untar extracts |filename| in |destdir|.
func (cmake *CMake) Untar(filename, destdir string) {
	cmake.Unzip(filename, destdir)
}

// Copy copies source to dest.
func (cmake *CMake) Copy(source, dest string) {
	cmake.WriteLine(fmt.Sprintf("message(STATUS \"Copy: %s %s\")", source, dest))
	cmake.WriteLine(fmt.Sprintf("execute_process(COMMAND"))
	cmake.WriteLine(fmt.Sprintf(
		"  ${CMAKE_COMMAND} -E copy \"%s\" \"%s\"", source, dest,
	))
	cmake.WriteLine(fmt.Sprintf("  RESULT_VARIABLE FAILURE)"))
	cmake.checkCommandError()
	cmake.WriteEmptyLine()
}

// CopyDir copies source to dest.
func (cmake *CMake) CopyDir(source, dest string) {
	cmake.WriteLine(fmt.Sprintf(
		"message(STATUS \"CopyDir: %s %s\")", source, dest,
	))
	cmake.WriteLine(fmt.Sprintf("execute_process(COMMAND"))
	cmake.WriteLine(fmt.Sprintf(
		"  ${CMAKE_COMMAND} -E copy_directory \"%s\" \"%s\"", source, dest,
	))
	cmake.WriteLine(fmt.Sprintf("  RESULT_VARIABLE FAILURE)"))
	cmake.checkCommandError()
	cmake.WriteEmptyLine()
}

// AddDefinition adds |definition| to the macro definitions
func (cmake *CMake) AddDefinition(definition string) {
	cmake.WriteLine(fmt.Sprintf(
		"LIST(APPEND CMAKE_REQUIRED_DEFINITIONS %s)", definition,
	))
}

// AddIncludeDir adds |path| to the header search path
func (cmake *CMake) AddIncludeDir(path string) {
	cmake.WriteLine(fmt.Sprintf(
		"LIST(APPEND CMAKE_REQUIRED_INCLUDES \"%s\")", path,
	))
}

// AddLibrary adds |library| to the libraries to include
func (cmake *CMake) AddLibrary(library string) {
	cmake.WriteLine(fmt.Sprintf(
		"LIST(APPEND CMAKE_REQUIRED_LIBRARIES \"%s\")", library,
	))
}

// checkPlatformCheckResult writes code to deal with a platform check result.
func (cmake *CMake) checkPlatformCheckResult(item, variable string, mandatory bool) {
	if mandatory {
		cmake.WriteLine(fmt.Sprintf("if(NOT (\"${%s}\"))", variable))
		cmake.WriteLine(fmt.Sprintf(
			"  message(FATAL_ERROR \"cannot find: %s\")", item,
		))
		cmake.WriteLine(fmt.Sprintf("endif()"))
	}
}

// CheckHeaderExists checks whether |header| exists and stores the
// result into the specified |variable|. If |mandatory| then, the
// processing will stop on failure. Otherwise, if found, then we'll
// add a preprocessor symbol named after |variable|.
func (cmake *CMake) CheckHeaderExists(header, variable string, mandatory bool) {
	cmake.WriteLine(fmt.Sprintf(
		"CHECK_INCLUDE_FILE_CXX(\"%s\" %s)", header, variable,
	))
	cmake.checkPlatformCheckResult(header, variable, mandatory)
}

// CheckLibraryExists checks whether |library| exists by looking for
// a function named |function|, storing the result in |variable|.
func (cmake *CMake) CheckLibraryExists(library, function, variable string, mandatory bool) {
	cmake.WriteLine(fmt.Sprintf(
		"CHECK_LIBRARY_EXISTS(\"%s\" \"%s\" \"\" %s)", library, function, variable,
	))
	cmake.checkPlatformCheckResult(library, variable, mandatory)
}

// SetRestrictiveCompilerFlags sets restrictive compiler flags.
func (cmake *CMake) SetRestrictiveCompilerFlags() {
	cmake.WriteSectionComment("Set restrictive compiler flags")
	cmake.output.WriteString(restrictiveflags.S)
	cmake.WriteEmptyLine()
	cmake.WriteLine(fmt.Sprintf("MkSetCompilerFlags()"))
}

// PrepareForCompilingTargets prepares internal variables such that
// we can compile targets with the required compiler flags.
func (cmake *CMake) PrepareForCompilingTargets() {
	cmake.WriteSectionComment("Prepare for compiling targets")
	cmake.WriteLine("add_definitions(${CMAKE_REQUIRED_DEFINITIONS})")
	cmake.WriteLine("include_directories(${CMAKE_REQUIRED_INCLUDES})")
	cmake.WriteLine("link_libraries(${CMAKE_REQUIRED_LIBRARIES})")
}

// AddExecutable defines an executable to be compiled.
func (cmake *CMake) AddExecutable(name string, sources []string) {
	cmake.WriteSectionComment(name)
	cmake.WriteLine(fmt.Sprintf("add_executable("))
	cmake.WriteLine(fmt.Sprintf("  %s", name))
	for _, source := range sources {
		cmake.WriteLine(fmt.Sprintf("  %s", source))
	}
	cmake.WriteLine(fmt.Sprintf(")"))
}

// AddTest defines a test to be run
func (cmake *CMake) AddTest(name string, arguments []string) {
	cmake.WriteSectionComment("test: "+name)
	cmake.WriteLine(fmt.Sprintf("add_test("))
	cmake.WriteLine(fmt.Sprintf("  NAME %s COMMAND", name))
	for _, arg := range arguments {
		cmake.WriteLine(fmt.Sprintf("  %s", arg))
	}
	cmake.WriteLine(fmt.Sprintf(")"))
}

// AddSingleHeaderDependency adds a single-header dependency
func (cmake *CMake) AddSingleHeaderDependency(SHA256, URL string) {
	headerName := filepath.Base(URL)
	cmake.WriteSectionComment(headerName)
	dirname := "${CMAKE_BINARY_DIR}/.mkbuild/include"
	filename := dirname + "/" + headerName
	cmake.MkdirAll(dirname)
	cmake.Download(filename, SHA256, URL)
	cmake.AddIncludeDir(dirname)
	guardVariable := "MK_HAVE_" + strings.ToUpper(strings.Replace(headerName, ".", "_", -1))
	cmake.CheckHeaderExists(headerName, guardVariable, true)
}

// AddSingleFileAsset adds a single-file asset to the build
func (cmake *CMake) AddSingleFileAsset(SHA256, URL string) {
	assetName := filepath.Base(URL)
	cmake.WriteSectionComment(assetName)
	dirname := "${CMAKE_BINARY_DIR}/.mkbuild/data"
	filename := dirname + "/" + assetName
	cmake.MkdirAll(dirname)
	cmake.Download(filename, SHA256, URL)
}

// IfWin32 allows you to generate WIN32 / !WIN32 specific code.
func (cmake *CMake) IfWIN32(thenFunc func(), elseFunc func()) {
	cmake.WriteLine("if((\"${WIN32}\"))")
	cmake.WithIndent("  ", thenFunc)
	cmake.WriteLine("else()")
	cmake.WithIndent("  ", elseFunc)
	cmake.WriteLine("endif()")
}

// If32bit allows you to generate 32 bit / 64 bit specific code. This
// function will configure cmake to fail if the bitsize is neither
// 32 not 64. That would be a very weird configuraton.
func (cmake *CMake) If32bit(func32 func(), func64 func()) {
	cmake.WriteLine("if((\"${CMAKE_SIZEOF_VOID_P}\" EQUAL 4))")
	cmake.WithIndent("  ", func32)
	cmake.WriteLine("elseif((\"${CMAKE_SIZEOF_VOID_P}\" EQUAL 8))")
	cmake.WithIndent("  ", func64)
	cmake.WriteLine("else()")
	cmake.WithIndent("  ", func() {
		cmake.WriteLine("message(FATAL_ERROR \"Neither 32 not 64 bit\")")
	})
	cmake.WriteLine("endif()")
}

// Close writes CMakeLists.txt in the current directory.
func (cmake *CMake) Close() {
	filename := "CMakeLists.txt"
	filep, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.WithError(err).Fatalf("os.Open failed for: %s", filename)
	}
	defer filep.Close()
	_, err = filep.WriteString(cmake.output.String())
	if err != nil {
		log.WithError(err).Fatalf("filep.WriteString failed for: %s", filename)
	}
	log.Infof("Written %s", filename)
}
