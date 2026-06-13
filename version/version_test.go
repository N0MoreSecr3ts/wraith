package version

import "testing"

// TestAppVersion ensures the composed version string matches the individual
// major, minor, patch, pre-release and build components.
func TestAppVersion(t *testing.T) {
	got := AppVersion()
	want := AppVersionMajor + "." + AppVersionMinor + "." + AppVersionPatch + AppVersionPre + AppVersionBuild

	if got != want {
		t.Fatalf("AppVersion() = %q, want %q", got, want)
	}
}
