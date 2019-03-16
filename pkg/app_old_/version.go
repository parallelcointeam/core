package app_old

// Name is the name of the application
const Name = "pod"

// Stamp is the version number placed by
var Stamp string

// Version returns the application version as a properly formed string per the semantic versioning 2.0.0 spec (http://semver.org/).
func Version() string {

	return Name + "-" + Stamp
}