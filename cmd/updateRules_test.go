package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCleanInputValid(t *testing.T) {
	Convey("Given valid git URLs", t, func() {
		cases := []struct {
			name string
			url  string
		}{
			{
				name: "a normal-length git URL",
				url:  "https://github.com/N0MoreSecr3ts/wraith-signatures",
			},
			{
				name: "a URL at the maximum allowed length",
				url:  "https://github.com/" + strings.Repeat("a", 2048-len("https://github.com/")),
			},
		}

		for _, c := range cases {
			Convey("When cleanInput is called with "+c.name, func() {
				result := cleanInput(c.url)

				Convey("It should return the URL unchanged", func() {
					So(result, ShouldEqual, c.url)
					So(len(c.url), ShouldBeLessThanOrEqualTo, 2048)
				})
			})
		}
	})
}

// TestCleanInputInvalid verifies that cleanInput exits the process for
// invalid input. cleanInput calls os.Exit, so each case is exercised by
// re-executing this test binary in a subprocess.
func TestCleanInputInvalid(t *testing.T) {
	if os.Getenv("CLEANINPUT_CRASH") == "1" {
		cleanInput(os.Getenv("CLEANINPUT_CRASH_URL"))
		return
	}

	Convey("Given invalid git URLs", t, func() {
		cases := []struct {
			name string
			url  string
		}{
			{name: "an empty URL", url: ""},
			{name: "a URL longer than 2048 characters", url: "https://github.com/" + strings.Repeat("a", 2049)},
		}

		for _, c := range cases {
			Convey("When cleanInput is called with "+c.name, func() {
				//nolint:gosec // G204: os.Args[0] is this test binary, re-exec'd to exercise cleanInput's os.Exit(2) path
				cmd := exec.Command(os.Args[0], "-test.run=TestCleanInputInvalid")
				cmd.Env = append(os.Environ(), "CLEANINPUT_CRASH=1", "CLEANINPUT_CRASH_URL="+c.url)
				err := cmd.Run()

				Convey("It should exit with a non-zero status", func() {
					So(err, ShouldNotBeNil)
				})
			})
		}
	})
}
