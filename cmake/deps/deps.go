// Package deps contains dependency related rules.
package deps

import (
	"fmt"

	"github.com/bassosimone/mkbuild/cmake/cmakefile"
	"github.com/bassosimone/mkbuild/cmake/cmakefile/prebuilt"
)

// All contains all the dependencies that we know of.
var All = map[string]func(*cmakefile.CMakeFile){
	"curl.haxx.se/ca": func(cmake *cmakefile.CMakeFile) {
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
	"github.com/measurement-kit/mkmock": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"f07bc063a2e64484482f986501003e45ead653ea3f53fadbdb45c17a51d916d2",
			"https://raw.githubusercontent.com/measurement-kit/mkmock/v0.2.0/mkmock.hpp",
		)
	},
}
