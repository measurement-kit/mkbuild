package rules

import (
	"fmt"
	"path/filepath"

	"github.com/bassosimone/mkbuild/cmake"
)

// WriteSectionComment writes a comment for |name| in |cmake|.
func WriteSectionComment(cmake *cmake.CMake, name string) {
	cmake.WriteLine("")
	cmake.WriteLine(fmt.Sprintf("#"))
	cmake.WriteLine(fmt.Sprintf("# %s", name))
	cmake.WriteLine(fmt.Sprintf("#"))
	cmake.WriteLine("")
}

// downloadSingleHeader downloads a library consisting of a single header.
func downloadSingleHeader(cmake *cmake.CMake, headerName, guardVariable, SHA256, URL string) {
	WriteSectionComment(cmake, headerName)
	dirname := filepath.Join("${CMAKE_BINARY_DIR}", ".mkbuild", "include")
	filename := filepath.Join(dirname, headerName)
	cmake.MkdirAll(dirname)
	cmake.Download(filename, SHA256, URL)
	cmake.AddIncludeDir(dirname)
	cmake.CheckHeaderExists(headerName, guardVariable, true)
	cmake.WriteLine("")
}

var Rules = map[string]func(*cmake.CMake){
	"curl.haxx.se/ca": func(cmake *cmake.CMake) {
		WriteSectionComment(cmake, "ca-bundle.pem")
		dirname := filepath.Join("${CMAKE_BINARY_DIR}", ".mkbuild", "etc")
		filename := filepath.Join(dirname, "ca-bundle.pem")
		cmake.MkdirAll(dirname)
		cmake.Download(
			filename, "4d89992b90f3e177ab1d895c00e8cded6c9009bec9d56981ff4f0a59e9cc56d6",
			"https://curl.haxx.se/ca/cacert-2018-12-05.pem",
		)
	},
	"github.com/adishavit/argh": func(cmake *cmake.CMake) {
		downloadSingleHeader(cmake, "argh.h", "MK_HAVE_ARGH_H",
			"ddb7dfc18dcf90149735b76fb2cff101067453a1df1943a6911233cb7085980c",
			"https://raw.githubusercontent.com/adishavit/argh/v1.3.0/argh.h",
		)
	},
	"github.com/catchorg/catch2": func(cmake *cmake.CMake) {
		downloadSingleHeader(cmake, "catch.hpp", "MK_HAVE_CATCH_HPP",
			"5eb8532fd5ec0d28433eba8a749102fd1f98078c5ebf35ad607fb2455a000004",
			"https://github.com/catchorg/Catch2/releases/download/v2.3.0/catch.hpp",
		)
	},
	"github.com/curl/curl": func(cmake *cmake.CMake) {
		WriteSectionComment(cmake, "libcurl")
		cmake.CheckHeaderExists("curl/curl.h", "MK_HAVE_CURL_CURL_H", true)
		cmake.CheckLibraryExists("curl", "curl_easy_init", "MK_HAVE_LIBCURL", true)
		cmake.AddLibrary("-lcurl")
	},
	"github.com/measurement-kit/mkmock": func(cmake *cmake.CMake) {
		downloadSingleHeader(cmake, "mkmock.hpp", "MK_HAVE_MKMOCK_HPP",
			"f07bc063a2e64484482f986501003e45ead653ea3f53fadbdb45c17a51d916d2",
			"https://raw.githubusercontent.com/measurement-kit/mkmock/v0.2.0/mkmock.hpp",
		)
	},
}

// XXX: add definition .mkbuild/include

/*
// installWinCurl downloads and verifies CURL's Windows zipfile for
// |version|, |arch|, having |SHA256| as SHA256. This function will set
// the proper IncludeDirs, LinkLibs, etc. as a side effect.
func installWinCurl(version, arch, SHA256 string) {
	prefix := filepath.Join(".mkbuild", "dep", "github.com", "curl", "curl")
	name := fmt.Sprintf("curl-%s-%s-mingw", version, arch)
	zipFile := fmt.Sprintf("%s.zip", name)
	url := fmt.Sprintf("https://curl.haxx.se/windows/dl-%s/%s", version, zipFile)
	downloadVerifyAndUnzip(
		filepath.Join(prefix, zipFile), prefix, SHA256, url,
	)
	gModuleInfo.IncludeDirs = append(gModuleInfo.IncludeDirs,
		filepath.Join(prefix, name, "include"))
	gModuleInfo.LinkLibs = append(gModuleInfo.LinkLibs,
		filepath.Join(prefix, name, "lib", "libcurl.dll.a"))
}
*/

/*
// installGithubcomCurlCurl installs github.com/curl/curl
func installGithubcomCurlCurl(dep string) {
	log.Infof("install: %s", dep)
	// TODO(bassosimone): we should probably specify this property via
	// command line. For example, `mkbuild autogen` will use the current
	// runtime.GOOS and `mkbuild autogen win32` will use win32.
		if runtime.GOOS != "windows" {
			gModuleInfo.LinkLibs = append(gModuleInfo.LinkLibs, "-lcurl")
			return
		}
	log.Warn("CURL support for Windows is still broken")
	winCurlVersion := "7.63.0"
	SHA256All := map[string]string{
		"win32": "9bf0f3a4d6aab8d3db7af3ed6edef9c3b12022b36c32edaf9ac443caa8899f65",
		"win64": "8795a1786a89607d0c52e3c0d8636aa29e4cf8c5b22a1dcb14ce6af0829b4814",
	}
	// TODO(bassosimone): the following is broken because we're setting
	// both the 32 bit and the 64 bit libraries. However, if we'll use
	// the command line, as mentioned above, the easy fix is to just use
	// whatever is provided from command line to choose.
	for arch, SHA256 := range SHA256All {
		installWinCurl(winCurlVersion, arch, SHA256)
	}
}
*/
