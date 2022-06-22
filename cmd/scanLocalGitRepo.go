// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"time"

	"github.com/N0MoreSecr3ts/wraith/core"

	"github.com/spf13/cobra"
)

// scanLocalGitRepoCmd represents the scanLocalGitRepo command
var scanLocalGitRepoCmd = &cobra.Command{
	Use:   "scanLocalGitRepo",
	Short: "Scan a git repo on a local machine",
	Long:  "Scan a git repo on a local machine",
	Run: func(cmd *cobra.Command, args []string) {

		scanType := "localGit"
		sess := core.NewSession(wraithConfig, scanType)

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

		core.GatherLocalRepositories(sess)
		core.AnalyzeRepositories(sess)
		sess.Finish()

		core.SummaryOutput(sess)

		if !sess.Silent && sess.WebServer {
			sess.Out.Important("Press Ctrl+C to stop web server and exit.\n")
			select {}
		}

	},
}

func init() {
	rootCmd.AddCommand(scanLocalGitRepoCmd)

	scanLocalGitRepoCmd.Flags().StringSlice("local-repos", nil, "List of local git repos to scan")
	scanLocalGitRepoCmd.Flags().Float64("commit-depth", -1, "Set the commit depth to scan")

	err := wraithConfig.BindPFlag("local-repos", scanLocalGitRepoCmd.Flags().Lookup("local-repos"))
	err = wraithConfig.BindPFlag("commit-depth", scanLocalGitRepoCmd.Flags().Lookup("commit-depth"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
