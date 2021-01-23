// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
	"wraith/core"

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

		// exclude the .git directory from local scans as it is not handled properly here
		sess.SkippablePath = core.AppendIfMissing(sess.SkippablePath, ".git/")

		if sess.Debug {
			core.PrintDebug(sess)
		}

		// By default we display a header to the user giving basic info about application. This will not be displayed
		// during a silent run which is the default when using this in an automated fashion.
		if !sess.JSONOutput && !sess.CSVOutput {
			sess.Out.Warn("%s\n\n", core.ASCIIBanner)
			sess.Out.Important("%s v%s started at %s\n", core.Name, sess.WraithVersion, sess.Stats.StartedAt.Format(time.RFC3339))
			sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
			if sess.WebServer {
				sess.Out.Important("Web interface available at http://%s:%d\n", sess.BindAddress, sess.BindPort)
			}
		}

		for _, p := range sess.LocalPaths {
			if core.PathExists(p, sess) {
				last := p[len(p)-1:]
				if last == "/" {
					core.ScanDir(p, sess)
				} else {
					core.DoFileScan(p, sess)
				}
			}
		}

		sess.Finish()

		core.SummaryOutput(sess)

		if !sess.Silent && sess.WebServer {
			sess.Out.Important("Press Ctrl+C to stop web server and exit.\n")
			select {}
		}

	},
}

func init() {
	rootCmd.AddCommand(scanLocalPathCmd)

	viperScanLocalPath = core.SetConfig()

	scanLocalPathCmd.Flags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	scanLocalPathCmd.Flags().Int("bind-port", 9393, "The port for the webserver")
	scanLocalPathCmd.Flags().Int("confidence-level", 3, "The confidence level level of the expressions used to find matches")
	scanLocalPathCmd.Flags().Bool("debug", false, "Print available debugging information to stdout")
	scanLocalPathCmd.Flags().Bool("hide-secrets", false, "Do not print secrets to any supported output")
	scanLocalPathCmd.Flags().StringSlice("ignore-extension", nil, "List of file extensions to ignore")
	scanLocalPathCmd.Flags().StringSlice("ignore-path", nil, "List of file paths to ignore")
	scanLocalPathCmd.Flags().Int("max-file-size", 10, "Max file size to scan (in MB)")
	scanLocalPathCmd.Flags().Int("num-threads", -1, "Number of execution threads")
	scanLocalPathCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanLocalPathCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default.yaml", "file(s) containing detection signatures.")
	scanLocalPathCmd.Flags().String("signature-path", "$HOME/.wraith/signatures", "path containing detection signatures.")
	scanLocalPathCmd.Flags().Bool("silent", false, "Suppress all output. An alternative output will need to be configured")
	scanLocalPathCmd.Flags().StringSlice("local-paths", nil, "List of local paths to scan")
	scanLocalPathCmd.Flags().Bool("web-server", false, "Enable the web interface for scan output")
	scanLocalPathCmd.Flags().Bool("csv", false, "output csv format")
	scanLocalPathCmd.Flags().Bool("json", false, "output json format")

	err := viperScanLocalPath.BindPFlag("bind-address", scanLocalPathCmd.Flags().Lookup("bind-address"))
	err = viperScanLocalPath.BindPFlag("bind-port", scanLocalPathCmd.Flags().Lookup("bind-port"))
	err = viperScanLocalPath.BindPFlag("debug", scanLocalPathCmd.Flags().Lookup("debug"))
	err = viperScanLocalPath.BindPFlag("hide-secrets", scanLocalPathCmd.Flags().Lookup("hide-secrets"))
	err = viperScanLocalPath.BindPFlag("ignore-extension", scanLocalPathCmd.Flags().Lookup("ignore-extension"))
	err = viperScanLocalPath.BindPFlag("ignore-path", scanLocalPathCmd.Flags().Lookup("ignore-path"))
	err = viperScanLocalPath.BindPFlag("confidence-level", scanLocalPathCmd.Flags().Lookup("confidence-level"))
	err = viperScanLocalPath.BindPFlag("max-file-size", scanLocalPathCmd.Flags().Lookup("max-file-size"))
	err = viperScanLocalPath.BindPFlag("num-threads", scanLocalPathCmd.Flags().Lookup("num-threads"))
	err = viperScanLocalPath.BindPFlag("scan-tests", scanLocalPathCmd.Flags().Lookup("scan-tests"))
	err = viperScanLocalPath.BindPFlag("signature-file", scanLocalPathCmd.Flags().Lookup("signature-file"))
	err = viperScanLocalPath.BindPFlag("signature-path", scanLocalPathCmd.Flags().Lookup("signature-path"))
	err = viperScanLocalPath.BindPFlag("local-paths", scanLocalPathCmd.Flags().Lookup("local-paths"))
	err = viperScanLocalPath.BindPFlag("silent", scanLocalPathCmd.Flags().Lookup("silent"))
	err = viperScanLocalPath.BindPFlag("web-server", scanLocalPathCmd.Flags().Lookup("web-server"))
	err = viperScanLocalPath.BindPFlag("csv", scanLocalPathCmd.Flags().Lookup("csv"))
	err = viperScanLocalPath.BindPFlag("json", scanLocalPathCmd.Flags().Lookup("json"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
