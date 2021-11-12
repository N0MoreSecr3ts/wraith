// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"time"

	"github.com/N0MoreSecr3ts/wraith/core"

	"github.com/spf13/cobra"
)

// scanLocalPathCmd represents the scanLocalFiles command
var scanLocalPathCmd = &cobra.Command{
	TraverseChildren: true,
	Use:              "scanLocalPath",
	Short:            "Scan local files and directorys",
	Long:             "Scan local files and directorys",
	Run: func(cmd *cobra.Command, args []string) {

		scanType := "localPath"
		sess := core.NewSession(wraithConfig, scanType)

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
		fmt.Println("Webserver: ", sess.WebServer)

		if !sess.Silent && sess.WebServer {
			sess.Out.Important("Press Ctrl+C to stop web server and exit.\n")
			select {}
		}

	},
}

func init() {
	rootCmd.AddCommand(scanLocalPathCmd)

	scanLocalPathCmd.Flags().StringSlice("local-paths", nil, "List of local paths to scan")

	err := wraithConfig.BindPFlag("local-paths", scanLocalPathCmd.Flags().Lookup("local-paths"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
