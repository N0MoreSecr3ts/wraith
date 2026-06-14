package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRootCmdUse(t *testing.T) {
	Convey("Given the root command", t, func() {
		Convey("Its Use string should be wraith", func() {
			So(rootCmd.Use, ShouldEqual, "wraith")
		})
	})
}

func TestRootCmdPersistentFlags(t *testing.T) {
	Convey("Given the root command's persistent flags", t, func() {
		cases := []struct {
			name         string
			defaultValue string
		}{
			{name: "bind-address", defaultValue: "127.0.0.1"},
			{name: "bind-port", defaultValue: "9393"},
			{name: "confidence-level", defaultValue: "3"},
			{name: "config-file", defaultValue: "$HOME/.wraith/config.yaml"},
			{name: "csv", defaultValue: "false"},
			{name: "debug", defaultValue: "false"},
			{name: "hide-secrets", defaultValue: "false"},
			{name: "ignore-extension", defaultValue: "[]"},
			{name: "ignore-path", defaultValue: "[]"},
			{name: "json", defaultValue: "false"},
			{name: "max-file-size", defaultValue: "10"},
			{name: "num-threads", defaultValue: "-1"},
			{name: "scan-tests", defaultValue: "false"},
			{name: "signature-file", defaultValue: "$HOME/.wraith/signatures/default.yaml"},
			{name: "signature-path", defaultValue: "$HOME/.wraith/signatures"},
			{name: "silent", defaultValue: "false"},
			{name: "web-server", defaultValue: "false"},
		}

		for _, c := range cases {
			Convey("The "+c.name+" flag should exist with its documented default", func() {
				flag := rootCmd.PersistentFlags().Lookup(c.name)
				So(flag, ShouldNotBeNil)
				So(flag.DefValue, ShouldEqual, c.defaultValue)
			})
		}
	})
}
