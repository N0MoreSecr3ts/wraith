package core_test

import (
	"testing"

	"github.com/N0MoreSecr3ts/wraith/core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFileExists(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given a filename", t, func() {

		Convey("When the file exists", func() {
			f := "../README.md"
			b := core.FileExists(f)

			Convey("The function should return true", func() {
				So(b, ShouldEqual, true)
			})
			Convey("The function should not return false", func() {
				So(b, ShouldNotEqual, false)
			})
		})

		Convey("When the file does not exist", func() {
			f := "../NOPE.md"
			b := core.FileExists(f)

			Convey("The function should return false", func() {
				So(b, ShouldEqual, false)
			})
			Convey("The function should not return true", func() {
				So(b, ShouldNotEqual, true)
			})
		})
	})
}
