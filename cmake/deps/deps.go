// Package deps contains dependency related rules.
package deps

import (
	"fmt"

	"github.com/apex/log"
	"github.com/measurement-kit/mkbuild/cmake/cmakefile"
	"github.com/measurement-kit/mkbuild/cmake/cmakefile/prebuilt"
)

// All contains all the dependencies that we know of.
var All = map[string]func(*cmakefile.CMakeFile){
	"curl.haxx.se/ca": func(cmake *cmakefile.CMakeFile) {
		log.Warn("curl.haxx.se/ca is deprecated; used github.com/measurement-kit/generic-assets instead")
		cmake.AddSingleFileAsset(
			"c1fd9b235896b1094ee97bfb7e042f93530b5e300781f59b45edf84ee8c75000",
			"https://curl.haxx.se/ca/cacert.pem",
		)
	},
	"github.com/adishavit/argh": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"ddb7dfc18dcf90149735b76fb2cff101067453a1df1943a6911233cb7085980c",
			"https://raw.githubusercontent.com/adishavit/argh/v1.3.0/argh.h",
		)
	},
	"github.com/catchorg/catch2": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"5eb8532fd5ec0d28433eba8a749102fd1f98078c5ebf35ad607fb2455a000004",
			"https://github.com/catchorg/Catch2/releases/download/v2.3.0/catch.hpp",
		)
	},
	"github.com/curl/curl": func(cmake *cmakefile.CMakeFile) {
		cmake.IfWIN32(func() {
			version := "7.61.1-1"
			cmake.Win32InstallPrebuilt(&prebuilt.Info{
				SHA256: "424d2f18f0f74dd6a0128f0f4e59860b7d2f00c80bbf24b2702e9cac661357cf",
				URL: fmt.Sprintf(
					"%s/%s/windows-curl-%s.tar.gz",
					"https://github.com/measurement-kit/prebuilt/releases/download/",
					"testing",
					version,
				),
				Prefix:     "MK_DIST/windows/curl/" + version,
				HeaderName: "curl/curl.h",
				LibName:    "libcurl.lib",
				FuncName:   "curl_easy_init",
			})
			cmake.AddRequiredDefinition("-DCURL_STATICLIB")
		}, func() {
			cmake.RequireHeaderExists("curl/curl.h")
			cmake.RequireLibraryExists("curl", "curl_easy_init")
			cmake.AddRequiredLibrary("curl")
		})
	},
	"github.com/howardhinnant/date": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"07aa75752540023ccccab178ed193f536c9d032cbbda997159af9f339d331eda",
			"https://raw.githubusercontent.com/HowardHinnant/date/v2.4.1/include/date/date.h",
		)
	},
	"github.com/maxmind/libmaxminddb": func(cmake *cmakefile.CMakeFile) {
		cmake.IfWIN32(func() {
			version := "1.3.2-2"
			cmake.Win32InstallPrebuilt(&prebuilt.Info{
				SHA256: "542933912814ac518037bd26083d0bba9daf68084f43c5cf2d7ec944d62b9ebb",
				URL: fmt.Sprintf(
					"%s/%s/windows-libmaxminddb-%s.tar.gz",
					"https://github.com/measurement-kit/prebuilt/releases/download/",
					"testing",
					version,
				),
				Prefix:     "MK_DIST/windows/libmaxminddb/" + version,
				HeaderName: "maxminddb.h",
				LibName:    "maxminddb.lib",
				FuncName:   "MMDB_open",
			})
		}, func() {
			cmake.RequireHeaderExists("maxminddb.h")
			cmake.RequireLibraryExists("maxminddb", "MMDB_open")
			cmake.AddRequiredLibrary("maxminddb")
		})
	},
	"github.com/measurement-kit/generic-assets": func(cmake *cmakefile.CMakeFile) {
		cmake.DownloadAndExtractArchive(
			"e7826c2575bacbc1aeccf64f10bfdf128c7ab38e6f5d17876775937986499df7",
			"https://github.com/measurement-kit/generic-assets/releases/download/20190205/generic-assets-20190205.tar.gz",
		)
	},
	"github.com/measurement-kit/mkcurl": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"1f301bb3d8e254a8ee40eac44373c2cfa9648d0d3ea96485e9231817aaf18076",
			"https://raw.githubusercontent.com/measurement-kit/mkcurl/v0.9.3/mkcurl.hpp",
		)
	},
	"github.com/measurement-kit/mkdata": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"96bb0384ecd7231a861111d8818a560b7d5ca83316cf7946a4f1a352db6ecfe3",
			"https://raw.githubusercontent.com/measurement-kit/mkdata/v0.3.0/mkdata.hpp",
		)
	},
	"github.com/measurement-kit/mkiplookup": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"a815119250d09be5eff332289f90fd872910f3dc9f29bb4a5fe60e272b38174f",
			"https://raw.githubusercontent.com/measurement-kit/mkiplookup/v0.2.0/mkiplookup.hpp",
		)
	},
	"github.com/measurement-kit/mkmmdb": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"c1cdcf2980c977a0d4abbdd447ddc19eefdfe6faa42b3be752d50f29930d4a87",
			"https://raw.githubusercontent.com/measurement-kit/mkmmdb/v0.4.0/mkmmdb.hpp",
		)
	},
	"github.com/measurement-kit/mkmock": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"f07bc063a2e64484482f986501003e45ead653ea3f53fadbdb45c17a51d916d2",
			"https://raw.githubusercontent.com/measurement-kit/mkmock/v0.2.0/mkmock.hpp",
		)
	},
	"github.com/nlohmann/json": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"8a6dbf3bf01156f438d0ca7e78c2971bca50eec4ca6f0cf59adf3464c43bb9d5",
			"https://raw.githubusercontent.com/nlohmann/json/v3.5.0/single_include/nlohmann/json.hpp",
		)
	},
}
