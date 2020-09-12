// Package common contains functionality not critical to the core project but still essential.
package core

import (
	"wraith/version"
)

// Project name and banner
const (
	Name = "wraith"
)

//	//ASCIIBanner = "        _ __           __\n" +
//	//	"  ___ _(_) /________  / /\n" +
//	//	" / _ `/ / __/ __/ _ \\/ _ \\\n" +
//	//	" \\_, /_/\\__/_/  \\___/_.__/\n" +
//	//	"/___/"
//
////    ASCIIBanner = "        _ __           __\n" +
////    	" _    _              _  _    _\n" +
////"| |  | |            (_)| |  | |\n" +
////"| |  | | _ __  __ _  _ | |_ | |__\n" +
////"| |/\| || '__|/ _` || || __|| '_ \\n" +
////"\  /\  /| |  | (_| || || |_ | | | |\n" +
////" \/  \/ |_|   \__,_||_| \__||_| |_|\n" +
////)
//
//var ASCIIBanner = `
// _    _              _  _    _
//| |  | |            (_)| |  | |
//| |  | | _ __  __ _  _ | |_ | |__
//| |//\| || '__|/ _` || || __|| '_ \
//\  /\  /| |  | (_| || || |_ | | | |
//\/  \/ |_|   \__,_||_| \__||_| |_|
//`

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
