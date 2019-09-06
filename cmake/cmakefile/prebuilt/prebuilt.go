// Package prebuilt helps to define prebuilt packages
package prebuilt

// Library contains info on a prebuilt library
type Library struct {
	// Name is the name of the library to try to link with to make
	// sure that the installation was correct
	Name string

	// Func is the name of the function to import from LibName
	// when trying to make sure that the installation was OK
	Func string
}

// Package describes a prebuilt package
type Package struct {
	// SHA256 is the SHA256 of the tarball
	SHA256 string

	// URL is the URL of the tarball
	URL string

	// Prefix is the prefix to strip from the tarball to reach the
	// arch dependent directories x86 and x64
	Prefix string

	// HeaderName is the name of a header to check for to make sure
	// that the installation was correct
	HeaderName string

	// Libs contains info to the libs to add to the build
	Libs []Library
}
