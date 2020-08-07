// Package common contains functionality not critical to the core project but still essential.
package common

import (
	"gitrob/version"
)

// TODO refactor out the common package

// Project name and banner
const (
	Name        = "gitrob"
	ASCIIBanner = "        _ __           __\n" +
		"  ___ _(_) /________  / /\n" +
		" / _ `/ / __/ __/ _ \\/ _ \\\n" +
		" \\_, /_/\\__/_/  \\___/_.__/\n" +
		"/___/"
)

// Version is the current version of gitlab
var Version = version.AppVersion()

// GitLabTanuki is the Gitlab specific banner
const GitLabTanuki = "\n" +
	"      //               //     \n" +
	"     ////             ////    \n" +
	"    //////           //////   \n" +
	"   ((((((((/////////((((((((  \n" +
	"   ((((((((////////(((((((((  \n" +
	"  ((((((((((///////(((((((((( \n" +
	"     ((((((((/////((((((((    \n" +
	"         (((((///(((((        \n" +
	"            (((/(((           \n" +
	"               *              \n" +
	"        GitLab Red Team       \n" +
	"\n"
