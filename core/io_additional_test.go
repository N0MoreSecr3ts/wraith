package core_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/N0MoreSecr3ts/wraith/core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAppendIfMissing(t *testing.T) {
	Convey("Given a slice of strings", t, func() {
		initial := []string{"a", "b"}

		Convey("When the value is already present", func() {
			result := core.AppendIfMissing(initial, "a")

			Convey("The slice should be unchanged", func() {
				So(result, ShouldResemble, []string{"a", "b"})
			})
		})

		Convey("When the value is not present", func() {
			result := core.AppendIfMissing(initial, "c")

			Convey("The value should be appended", func() {
				So(result, ShouldResemble, []string{"a", "b", "c"})
			})
		})
	})
}

func TestIsMaxFileSize(t *testing.T) {
	Convey("Given a temporary file and session", t, func() {
		dir := t.TempDir()
		file := filepath.Join(dir, "small.txt")

		// create a small file
		err := os.WriteFile(file, []byte("small"), 0o644)
		So(err, ShouldBeNil)

		sess := &core.Session{
			MaxFileSize: 1, // 1 MB
		}

		Convey("When the file is smaller than the max size", func() {
			isTooLarge, reason := core.IsMaxFileSize(file, sess)

			Convey("The file should not be considered too large", func() {
				So(isTooLarge, ShouldBeFalse)
				So(reason, ShouldEqual, "")
			})
		})

		Convey("When the file does not exist", func() {
			isTooLarge, reason := core.IsMaxFileSize(filepath.Join(dir, "missing.txt"), sess)

			Convey("The function should indicate the file does not exist", func() {
				So(isTooLarge, ShouldBeFalse)
				So(reason, ShouldEqual, "does not exist")
			})
		})
	})
}
