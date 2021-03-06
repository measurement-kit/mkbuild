// Package deps contains dependency related rules.
package deps

import (
	"github.com/apex/log"
	"github.com/measurement-kit/mkbuild/cmake/cmakefile"
)

// All contains all the dependencies that we know of.
var All = map[string]func(*cmakefile.CMakeFile){
	"github.com/adishavit/argh": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"ddb7dfc18dcf90149735b76fb2cff101067453a1df1943a6911233cb7085980c",
			"https://raw.githubusercontent.com/adishavit/argh/v1.3.1/argh.h",
		)
	},
	"github.com/c-ares/c-ares": func(cmake *cmakefile.CMakeFile) {
		// TODO(bassosimone): implement c-ares support for Windows
		log.Warn("github.com/c-ares/c-ares: not supported on Windows")
		cmake.RequireHeaderExists("ares.h")
		cmake.RequireLibraryExists("cares", "ares_process")
		cmake.AddRequiredLibrary("cares")
	},
	"github.com/catchorg/catch2": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"914690be7714fed1f228679a0e379b4e1f6d7d66f88ac20335b7853ff1a8bc55",
			"https://github.com/catchorg/Catch2/releases/download/v2.11.1/catch.hpp",
		)
	},
	"github.com/curl/curl": func(cmake *cmakefile.CMakeFile) {
		// TODO(bassosimone): implement curl support for Windows
		cmake.RequireHeaderExists("curl/curl.h")
		cmake.RequireLibraryExists("curl", "curl_easy_init")
		cmake.AddRequiredLibrary("curl")
	},
	"github.com/howardhinnant/date": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"07aa75752540023ccccab178ed193f536c9d032cbbda997159af9f339d331eda",
			"https://raw.githubusercontent.com/HowardHinnant/date/v2.4.1/include/date/date.h",
		)
	},
	"github.com/maxmind/libmaxminddb": func(cmake *cmakefile.CMakeFile) {
		// TODO(bassosimone): implement libmaxminddb support for Windows
		cmake.RequireHeaderExists("maxminddb.h")
		cmake.RequireLibraryExists("maxminddb", "MMDB_open")
		cmake.AddRequiredLibrary("maxminddb")
	},
	"github.com/measurement-kit/generic-assets": func(cmake *cmakefile.CMakeFile) {
		cmake.DownloadAndExtractArchive(
			"70d590c20b2ed31fd43cc63709b267672fecfeac7e908d11e845664ddd43b04f",
			"https://github.com/measurement-kit/generic-assets/releases/download/20190520205742/generic-assets-20190520205742.tar.gz",
		)
	},
	"github.com/measurement-kit/mkbouncer": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"7c0a810d58fdbc1ffdacb0eed442d31ebff819e49d1222b81619fe3c582a028c",
			"https://raw.githubusercontent.com/measurement-kit/mkbouncer/v0.3.0/mkbouncer.hpp",
		)
	},
	"github.com/measurement-kit/mkcollector": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"9c81a0c4212eb411be380d2d4b0bd3ada1d70b23f6039b17fe82d3d4ccad1774",
			"https://raw.githubusercontent.com/measurement-kit/mkcollector/v0.6.0/mkcollector.hpp",
		)
	},
	"github.com/measurement-kit/mkcurl": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"cb58b85ccdf8d9f3b559ccab7d2f31cad4f8aba774e9953534e32373d415ec94",
			"https://raw.githubusercontent.com/measurement-kit/mkcurl/v0.12.0/mkcurl.hpp",
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
	"github.com/measurement-kit/mkuuid4": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"5b6b4445697d9beb6ad5310d98b7743c2ffe8266cdec79df0a7a429dcfc247ac",
			"https://raw.githubusercontent.com/measurement-kit/mkuuid4/v0.1.0/mkuuid4.hpp",
		)
	},
	"github.com/nlohmann/json": func(cmake *cmakefile.CMakeFile) {
		cmake.AddSingleHeaderDependency(
			"3b5d2b8f8282b80557091514d8ab97e27f9574336c804ee666fda673a9b59926",
			"https://raw.githubusercontent.com/nlohmann/json/v3.7.3/single_include/nlohmann/json.hpp",
		)
	},
	"github.com/openssl/openssl": func(cmake *cmakefile.CMakeFile) {
		// TODO(bassosimone): implement openssl support for Windows
		cmake.IfAPPLE(func() {
			// Automatically use Homebrew, if available
			cmake.WriteLine(`if(EXISTS "/usr/local/opt/openssl@1.1")`)
			cmake.WithIndent("  ", func() {
				cmake.WriteLine(`  set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -I/usr/local/opt/openssl@1.1/include")`)
				cmake.WriteLine(`  set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -I/usr/local/opt/openssl@1.1/include")`)
				cmake.WriteLine(`  set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -L/usr/local/opt/openssl@1.1/lib")`)
				cmake.WriteLine(`  set(CMAKE_SHARED_LINKER_FLAGS "${CMAKE_SHARED_LINKER_FLAGS} -L/usr/local/opt/openssl@1.1/lib")`)
			})
			cmake.WriteLine("endif()")
		}, nil)
		cmake.RequireHeaderExists("openssl/rsa.h")
		cmake.RequireLibraryExists("crypto", "RSA_new")
		cmake.AddRequiredLibrary("crypto")
		cmake.RequireHeaderExists("openssl/ssl.h")
		cmake.RequireLibraryExists("ssl", "SSL_read")
		cmake.AddRequiredLibrary("ssl")
	},
}
