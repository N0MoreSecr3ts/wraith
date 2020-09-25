// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
	"wraith/core"
	"wraith/version"

	"github.com/spf13/cobra"
)

var viperScanLocalPath *viper.Viper

// scanLocalPathCmd represents the scanLocalFiles command
var scanLocalPathCmd = &cobra.Command{
	Use:   "scanLocalPath",
	Short: "Scan local files and directorys",
	Long:  "Scan local files and directorys",
	Run: func(cmd *cobra.Command, args []string) {

		scanType := "localPath"
		sess := core.NewSession(viperScanLocalPath, scanType)

		core.CheckArgs(sess.LocalFiles, sess.LocalDirs, sess)

		// exclude the .git directory from local scans as it is not handled properly here
		sess.SkippablePath = core.AppendIfMissing(sess.SkippablePath, ".git/")

		//sess.Out.Info("%s\n\n", common.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", sess.BindAddress, sess.BindPort)

		// Run either a file scan directly, or if it is a directory then walk the path and gather eligible files and then run a scan against each of them
		for _, fl := range sess.LocalFiles {
			if fl != "" {
				if !core.PathExists(fl, sess) {
					sess.Out.Error("\n[*] <%s> does not exist! Quitting.\n", fl)
				} else {
					core.DoFileScan(fl, sess)
				}
			}
		}

		for _, pth := range sess.LocalDirs {
			if pth != "" {
				if !core.PathExists(pth, sess) {
					sess.Out.Error("\n[*] <%s> does not exist! Quitting.\n", pth)
				} else {
					core.ScanDir(pth, sess)
				}
			}
		}

		sess.Finish()

		if sess.JSONOutput || sess.CSVOutput {
			core.WriteOutput(sess)
		}

		core.PrintSessionStats(sess)

		if !sess.Silent {
			sess.Out.Important("Press Ctrl+C to stop web server and exit.")
			select {}
		}

	},
}

func init() {
	rootCmd.AddCommand(scanLocalPathCmd)

	viperScanLocalPath = core.SetConfig()

	scanLocalPathCmd.Flags().Bool("csv", false, "Write results to --output-file in CSV format")
	scanLocalPathCmd.Flags().Bool("debug", false, "Print debugging information")
	scanLocalPathCmd.Flags().Bool("hide-secrets", false, "Show secrets in any supported output")
	scanLocalPathCmd.Flags().String("ignore-extension", "", "a list of extensions to ignore during a scan")
	scanLocalPathCmd.Flags().String("ignore-path", "", "a list of paths to ignore during a scan")
	scanLocalPathCmd.Flags().Bool("json", false, "Write results to --output-file in JSON format")
	scanLocalPathCmd.Flags().Int64("max-file-size", 50, "Max file size to scan")
	scanLocalPathCmd.Flags().Int("match-level", 3, "The match level of the expressions used to find matches")
	scanLocalPathCmd.Flags().String("output-dir", "./", "Write csv and/or json files to directory")
	scanLocalPathCmd.Flags().String("output-prefix", "wraith", "Prefix to prepend to datetime stamp for output files")
	scanLocalPathCmd.Flags().String("scan-dir", "", "scan a directory of files not from a git project")
	scanLocalPathCmd.Flags().String("scan-file", "", "scan a single file")
	scanLocalPathCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanLocalPathCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default.yml", "file(s) containing secrets detection signatures.")
	scanLocalPathCmd.Flags().Bool("silent", false, "Suppress all output except for errors")


	err := viperScanLocalPath.BindPFlag("csv", scanLocalPathCmd.Flags().Lookup("csv"))
	err = viperScanLocalPath.BindPFlag("debug", scanLocalPathCmd.Flags().Lookup("debug"))
	err = viperScanLocalPath.BindPFlag("hide-secrets", scanLocalPathCmd.Flags().Lookup("hide-secrets"))
	err = viperScanLocalPath.BindPFlag("ignore-extension", scanLocalPathCmd.Flags().Lookup("ignore-extension"))
	err = viperScanLocalPath.BindPFlag("ignore-path", scanLocalPathCmd.Flags().Lookup("ignore-path"))
	err = viperScanLocalPath.BindPFlag("json", scanLocalPathCmd.Flags().Lookup("json"))
	err = viperScanLocalPath.BindPFlag("max-file-size", scanLocalPathCmd.Flags().Lookup("max-file-size"))
	err = viperScanLocalPath.BindPFlag("match-level", scanLocalPathCmd.Flags().Lookup("match-level"))
	err = viperScanLocalPath.BindPFlag("output-dir", scanGithubCmd.Flags().Lookup("output-dir"))
	err = viperScanLocalPath.BindPFlag("output-prefix", scanGithubCmd.Flags().Lookup("output-prefix"))
	err = viperScanLocalPath.BindPFlag("scan-dir", scanLocalPathCmd.Flags().Lookup("scan-dir"))
	err = viperScanLocalPath.BindPFlag("scan-file", scanLocalPathCmd.Flags().Lookup("scan-file"))
	err = viperScanLocalPath.BindPFlag("scan-tests", scanLocalPathCmd.Flags().Lookup("scan-tests"))
	err = viperScanLocalPath.BindPFlag("signature-file", scanLocalPathCmd.Flags().Lookup("signature-file"))
	err = viperScanLocalPath.BindPFlag("silent", scanLocalPathCmd.Flags().Lookup("silent"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
