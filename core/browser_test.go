package core_test

import (
	"testing"

	"github.com/N0MoreSecr3ts/wraith/core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCleanURLSpaces(t *testing.T) {

	Convey("Given a string", t, func() {

		Convey("When the strings have spaces", func() {
			str := "This "
			result := core.CleanURLSpaces(str)

			Convey("The spaces should be replaced with dashes", func() {
				for i, _ := range result {
					So(result[i], ShouldEqual, "This-")
				}
			})
			Convey("The spaces should not be replaced with underscores", func() {
				for i, _ := range result {
					So(result[i], ShouldNotEqual, "This_")
				}
			})
			Convey("The spaces should not be left alone", func() {
				for i, _ := range result {
					So(result[i], ShouldNotEqual, "This ")
				}
			})
			Convey("The spaces should not be replaced with \"&#160;\"", func() {
				for i, _ := range result {
					So(result[i], ShouldNotEqual, "This&#160;")
				}
			})
		})

		Convey("When the strings do not have spaces", func() {
			str := "This"
			result := core.CleanURLSpaces(str)

			Convey("The spaces should be left alone", func() {
				for i, _ := range result {
					So(result[i], ShouldEqual, "This")
				}
			})

			Convey("The string not should contain an extra dash", func() {
				for i, _ := range result {
					So(result[i], ShouldNotEqual, "This-")
				}
			})
			Convey("The string should not contain an extra underscores", func() {
				for i, _ := range result {
					So(result[i], ShouldNotEqual, "This_")
				}
			})
			Convey("The string should not contain an extra \"&#160;\"", func() {
				for i, _ := range result {
					So(result[i], ShouldNotEqual, "This&#160;")
				}
			})
		})
	})
}
