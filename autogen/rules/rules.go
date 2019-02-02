// Package rules contains the build rules.
package rules

import (
	"fmt"

	"github.com/bassosimone/mkbuild/autogen/cmake"
)

// downloadWinCurl downloads curl for Windows
func downloadWinCurl(cmake *cmake.CMake, filename, SHA256, URL string) {
	dirname := "${CMAKE_BINARY_DIR}/.mkbuild/download"
	filepathname := dirname + "/" + filename
	cmake.MkdirAll(dirname)
	cmake.Download(filepathname, SHA256, URL)
	cmake.Untar(filepathname, dirname)
}

// Rules contains all the build rules that we know of.
var Rules = map[string]func(*cmake.CMake){
	"curl.haxx.se/ca": func(cmake *cmake.CMake) {
		cmake.AddSingleFileAsset(
			"4d89992b90f3e177ab1d895c00e8cded6c9009bec9d56981ff4f0a59e9cc56d6",
			"https://curl.haxx.se/ca/cacert.pem",
		)
	},
	"github.com/adishavit/argh": func(cmake *cmake.CMake) {
		cmake.AddSingleHeaderDependency(
			"ddb7dfc18dcf90149735b76fb2cff101067453a1df1943a6911233cb7085980c",
			"https://raw.githubusercontent.com/adishavit/argh/v1.3.0/argh.h",
		)
	},
	"github.com/catchorg/catch2": func(cmake *cmake.CMake) {
		cmake.AddSingleHeaderDependency(
			"5eb8532fd5ec0d28433eba8a749102fd1f98078c5ebf35ad607fb2455a000004",
			"https://github.com/catchorg/Catch2/releases/download/v2.3.0/catch.hpp",
		)
	},
	"github.com/curl/curl": func(cmake *cmake.CMake) {
		cmake.WriteSectionComment("libcurl")
		cmake.WriteLine("if((\"${WIN32}\"))")
		cmake.WithIndent("  ", func() {
			version := "7.61.1-1"
			release := "testing"
			baseURL := "https://github.com/measurement-kit/prebuilt/releases/download/"
			URL := fmt.Sprintf("%s/%s/windows-curl-%s.tar.gz", baseURL, release, version)
			downloadWinCurl(
				cmake, "windows-curl.tar.gz",
				"424d2f18f0f74dd6a0128f0f4e59860b7d2f00c80bbf24b2702e9cac661357cf",
				URL,
			)
			cmake.WriteLine("if((\"${CMAKE_SIZEOF_VOID_P}\" EQUAL 4))")
			cmake.WithIndent("  ", func() {
				cmake.WriteLine("SET(MK_CURL_ARCH \"x86\")")
			})
			cmake.WriteLine("else()")
			cmake.WithIndent("  ", func() {
				cmake.WriteLine("SET(MK_CURL_ARCH \"x64\")")
			})
			cmake.WriteLine("endif()")
			cmake.WriteEmptyLine()
			curldir := "${CMAKE_BINARY_DIR}/.mkbuild/download/MK_DIST/windows/curl/" + version + "/${MK_CURL_ARCH}"
			includedirname := curldir + "/include"
			libname := curldir + "/lib/libcurl.lib"
			cmake.AddIncludeDir(includedirname)
			cmake.CheckHeaderExists("curl/curl.h", "MK_HAVE_CURL_CURL_H", true)
			cmake.WriteEmptyLine()
			cmake.CheckLibraryExists(libname, "curl_easy_init", "MK_HAVE_LIBCURL", true)
			cmake.AddLibrary(libname)
			cmake.AddDefinition("-DCURL_STATICLIB")
		})
		cmake.WriteEmptyLine()
		cmake.WriteLine("else()")
		cmake.WithIndent(" ", func() {
			cmake.CheckHeaderExists("curl/curl.h", "MK_HAVE_CURL_CURL_H", true)
			cmake.CheckLibraryExists("curl", "curl_easy_init", "MK_HAVE_LIBCURL", true)
			cmake.AddLibrary("curl")
		})
		cmake.WriteLine("endif()")
	},
	"github.com/measurement-kit/mkmock": func(cmake *cmake.CMake) {
		cmake.AddSingleHeaderDependency(
			"f07bc063a2e64484482f986501003e45ead653ea3f53fadbdb45c17a51d916d2",
			"https://raw.githubusercontent.com/measurement-kit/mkmock/v0.2.0/mkmock.hpp",
		)
	},
}
