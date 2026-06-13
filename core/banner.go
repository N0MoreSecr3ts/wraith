// Package core represents the core functionality of all commands
package core

import (
	"github.com/N0MoreSecr3ts/wraith/version"
)

// Project name and banner
const (
	Name = "wraith"
)

// ASCIIBanner is the project specific banner
const ASCIIBanner = "\n" +
	"____    __    ____ .______          ___       __  .___________. __    __\n" +
	"\\   \\  /  \\  /   / |   _  \\        /   \\     |  | |           ||  |  |  |\n" +
	" \\   \\/    \\/   /  |  |_)  |      /  ^  \\    |  | `---|  |----`|  |__|  |\n" +
	"  \\            /   |      /      /  /_\\  \\   |  |     |  |     |   __   |\n" +
	"   \\    /\\    /    |  |\\  \\----./  _____  \\  |  |     |  |     |  |  |  |\n" +
	"    \\__/  \\__/     | _| `._____/__/     \\__\\ |__|     |__|     |__|  |__|\n" +
	"\n"

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
