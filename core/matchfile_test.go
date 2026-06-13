package core

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewMatchFile(t *testing.T) {
	Convey("Given a file path", t, func() {
		path := "/tmp/foo/bar.txt"

		Convey("When newMatchFile is called", func() {
			m := newMatchFile(path)

			Convey("The MatchFile should contain the original path", func() {
				So(m.Path, ShouldEqual, path)
			})

			Convey("The MatchFile should contain the filename", func() {
				So(m.Filename, ShouldEqual, "bar.txt")
			})

			Convey("The MatchFile should contain the extension", func() {
				So(m.Extension, ShouldEqual, ".txt")
			})
		})
	})
}

func TestMatchFileIsSkippable(t *testing.T) {
	Convey("Given a MatchFile and Session", t, func() {
		sess := &Session{
			SkippableExt:  []string{".log"},
			SkippablePath: []string{"node_modules", "vendor"},
		}

		Convey("When the file has a skippable extension", func() {
			m := newMatchFile("/tmp/app/output.log")

			Convey("The file should be skippable", func() {
				So(m.isSkippable(sess), ShouldBeTrue)
			})
		})

		Convey("When the file is inside a skippable path", func() {
			m := newMatchFile("/repo/node_modules/lib.js")

			Convey("The file should be skippable", func() {
				So(m.isSkippable(sess), ShouldBeTrue)
			})
		})

		Convey("When the file is not in a skippable path or extension list", func() {
			m := newMatchFile("/repo/src/main.go")

			Convey("The file should not be skippable", func() {
				So(m.isSkippable(sess), ShouldBeFalse)
			})
		})
	})
}
