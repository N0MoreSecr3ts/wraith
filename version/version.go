// Package version sets the version for the tool
package version

import (
	"fmt"
)

//AppVersionMajor is the major revision number
const AppVersionMajor = "0"

// AppVersionMinor is the minor revision number
const AppVersionMinor = "0"

// AppVersionPatch is the patch version
const AppVersionPatch = "7"

// AppVersionPre ...
const AppVersionPre = ""

// AppVersionBuild should be empty string when releasing
const AppVersionBuild = ""

// AppVersion generates a usable version string
func AppVersion() string {
	return fmt.Sprintf("%s.%s.%s%s%s", AppVersionMajor, AppVersionMinor, AppVersionPatch, AppVersionPre, AppVersionBuild)
}
