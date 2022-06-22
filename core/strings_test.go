package core_test

import (
	"testing"

	"github.com/N0MoreSecr3ts/wraith/core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPluralize(t *testing.T) {

	Convey("Given a string and a count", t, func() {

		Convey("When the word is 'test'", func() {
			w1 := "test"
			w2 := "tests"

			Convey("If the count is 0", func() {
				w := core.Pluralize(0, w1, w2)

				Convey("The word should be 'tests'", func() {
					So(w, ShouldEqual, w2)
				})
				Convey("The word should not be 'test'", func() {
					So(w, ShouldNotEqual, w1)
				})
			})

			Convey("If the count is 1", func() {
				w := core.Pluralize(1, w1, w2)

				Convey("The word should be 'test'", func() {
					So(w, ShouldEqual, w1)
				})
				Convey("The word should not be 'tests'", func() {
					So(w, ShouldNotEqual, w2)
				})
			})

			Convey("If the count is -1", func() {
				w := core.Pluralize(-1, w1, w2)

				Convey("The word should be 'tests'", func() {
					So(w, ShouldEqual, w2)
				})
				Convey("The word should not be 'test'", func() {
					So(w, ShouldNotEqual, w1)
				})
			})

			Convey("If the count is -0", func() {
				w := core.Pluralize(-0, w1, w2)

				Convey("The word should be 'tests'", func() {
					So(w, ShouldEqual, w2)
				})
				Convey("The word should not be 'test'", func() {
					So(w, ShouldNotEqual, w1)
				})
			})
		})
	})
}

func TestTruncateString(t *testing.T) {

	Convey("Given a string and a length ", t, func() {
		s := "This is too long"

		Convey("When the string is 'This is too long' and the max length is 10", func() {
			l := 10
			str := core.TruncateString(s, l)

			Convey("The new string should be 'This is to'", func() {
				So(str, ShouldEqual, "This is to...")
			})

			Convey("The new string should not be 'This is too long' ", func() {
				So(s, ShouldEqual, "This is too long")
			})
		})
	})
}
