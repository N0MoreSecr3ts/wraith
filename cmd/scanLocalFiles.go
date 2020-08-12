// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
	"wraith/core"
	"wraith/version"

	"github.com/spf13/cobra"
)

var viperScanLocalPath *viper.Viper

// scanLocalPathCmd represents the scanLocalFiles command
var scanLocalPathCmd = &cobra.Command{
	Use:   "scanLocalFiles",
	Short: "Scan local files and directorys",
	Long:  "Scan local files and directorys",
	Run: func(cmd *cobra.Command, args []string) {

		// Ensure that both the file and a directory flags are not set
		_ = core.CheckArgs(viperScanLocalPath.GetString("scan-dir"), viperScanLocalPath.GetString("scan-file"))

		splitDir := strings.Split(viperScanLocalPath.GetString("scan-dir"), ",")
		splitFile := strings.Split(viperScanLocalPath.GetString("scan-file"), ",")

		var sDir []string
		var sFile []string

		for _, pth := range splitDir {
			sDir = append(sDir, pth)
		}

		for _, fl := range splitFile {
			sDir = append(sDir, fl)
		}

		scanType := "localPath"
		sess, err := core.NewSession(viperScanLocalPath, scanType)

		// exclude the .git directory from local scans as it is not handled properly here
		sess.SkippablePath = core.AppendIfMissing(sess.SkippablePath, ".git/")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		//sess.Out.Info("%s\n\n", common.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", "127.0.0.1", 9393)

		//core.GatherLocalRepositories(sess)
		//core.AnalyzeRepositories(sess)

		// Run either a file scan directly, or if it is a directory then walk the path and gather eligible files and then run a scan against each of them
		for _, fl := range sFile {
			if fl != "" {
				if !core.PathExists(fl) {
					sess.Out.Error("\n[*] <%s> does not exist! Quitting.\n", fl)
					os.Exit(1)
				} else {
					core.DoFileScan(fl, sess)
				}
			}
		}

		for _, pth := range sDir {
			if pth != "" {
				if !core.PathExists(pth) {
					sess.Out.Error("\n[*] <%s> does not exist! Quitting.\n", pth)
					os.Exit(1)
				} else {
					core.ScanDir(pth, sess)
				}
			}
		}

		sess.Finish()

		core.PrintSessionStats(sess)

		if !sess.Silent {
			sess.Out.Important("Press Ctrl+C to stop web server and exit.")
			select {}
		}

	},
}

func init() {
	rootCmd.AddCommand(scanLocalPathCmd)

	//scanLocalPathCmd.Flags().Bool("csv", false, "output csv format")
	scanLocalPathCmd.Flags().Bool("debug", false, "Print debugging information")
	scanLocalPathCmd.Flags().Bool("hide-secrets", false, "Show secrets in any supported output")
	//scanLocalPathCmd.Flags().Bool("json", false, "output json format")
	//scanLocalPathCmd.Flags().Bool("load-triage", false, "load a triage file")
	scanLocalPathCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanLocalPathCmd.Flags().Bool("silent", false, "Suppress all output except for errors")
	scanLocalPathCmd.Flags().Int64("max-file-size", 50, "Max file size to scan")
	scanLocalPathCmd.Flags().Int("match-level", 3, "The match level of the expressions used to find matches")
	scanLocalPathCmd.Flags().String("ignore-extension", "", "a list of extensions to ignore during a scan")
	scanLocalPathCmd.Flags().String("ignore-path", "", "a list of paths to ignore during a scan")
	scanLocalPathCmd.Flags().String("rules-file", "$HOME/grover/rules/default.yml", "file(s) containing secrets detection rules.")
	scanLocalPathCmd.Flags().String("scan-dir", "", "scan a directory of files not from a git project")
	scanLocalPathCmd.Flags().String("scan-file", "", "scan a single file")
	//scanLocalPathCmd.Flags().String("triage-file", "$HOME/.grover/triage.yaml", "file containing secrets that have been previously triaged.")

	//viperScanLocalPath.BindPFlag("csv", scanLocalPathCmd.Flags().Lookup("csv"))
	viperScanLocalPath.BindPFlag("debug", scanLocalPathCmd.Flags().Lookup("debug"))
	viperScanLocalPath.BindPFlag("hide-secrets", scanLocalPathCmd.Flags().Lookup("hide-secrets"))
	//viperScanLocalPath.BindPFlag("json", scanLocalPathCmd.Flags().Lookup("json"))
	//viperScanLocalPath.BindPFlag("load-triage", scanLocalPathCmd.Flags().Lookup("load-triage"))
	viperScanLocalPath.BindPFlag("scan-tests", scanLocalPathCmd.Flags().Lookup("scan-tests"))
	viperScanLocalPath.BindPFlag("silent", scanLocalPathCmd.Flags().Lookup("silent"))
	viperScanLocalPath.BindPFlag("max-file-size", scanLocalPathCmd.Flags().Lookup("max-file-size"))
	viperScanLocalPath.BindPFlag("match-level", scanLocalPathCmd.Flags().Lookup("match-level"))
	viperScanLocalPath.BindPFlag("ignore-extension", scanLocalPathCmd.Flags().Lookup("ignore-extension"))
	viperScanLocalPath.BindPFlag("ignore-path", scanLocalPathCmd.Flags().Lookup("ignore-path"))
	viperScanLocalPath.BindPFlag("rules-file", scanLocalPathCmd.Flags().Lookup("rules-file"))
	viperScanLocalPath.BindPFlag("scan-dir", scanLocalPathCmd.Flags().Lookup("scan-dir"))
	viperScanLocalPath.BindPFlag("scan-file", scanLocalPathCmd.Flags().Lookup("scan-file"))
	//viperScanLocalPath.BindPFlag("triage-file", scanLocalPathCmd.Flags().Lookup("triage-file"))
}
