// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"time"

	"github.com/N0MoreSecr3ts/wraith/core"

	"github.com/spf13/cobra"
)

// scanGitlabCmd represents the scanGitlab command
var scanGitlabCmd = &cobra.Command{
	Use:   "scanGitlab",
	Short: "Scan one or more gitlab groups or users for secrets",
	Long:  `Scan one or more gitlab groups or users for secrets`,
	Run: func(cmd *cobra.Command, args []string) {

		scanType := "gitlab"
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

		if sess.Debug {
			sess.Out.Debug("We have these orgs: %s\n", sess.UserOrgs)
			sess.Out.Debug("We have these users: %s\n", sess.UserLogins)
			sess.Out.Debug("We have these repos: %s\n", sess.UserRepos)
		}

		sess.GitlabAccessToken = wraithConfig.GetString("gitlab-api-token")

		sess.InitGitClient()

		core.GatherTargets(sess)
		core.GatherGitlabRepositories(sess)
		core.AnalyzeRepositories(sess)
		sess.Finish()

		core.SummaryOutput(sess)

		if !sess.Silent && sess.WebServer {
			sess.Out.Important("%s", core.GitLabTanuki)
			sess.Out.Important("Press Ctrl+C to stop web server and exit.\n")
			select {}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanGitlabCmd)

	scanGitlabCmd.Flags().Bool("add-org-members", false, "Add members to targets when processing organizations")
	scanGitlabCmd.Flags().String("gitlab-api-token", "", "API token for access to gitlab, see doc for necessary scope")
	scanGitlabCmd.Flags().StringSlice("gitlab-projects", nil, "List of Gitlab projects or users to scan")
	scanGitlabCmd.Flags().Float64("commit-depth", -1, "Set the commit depth to scan")

	err := wraithConfig.BindPFlag("commit-depth", scanGitlabCmd.Flags().Lookup("commit-depth"))
	err = wraithConfig.BindPFlag("add-org-members", scanGitlabCmd.Flags().Lookup("add-org-members"))
	err = wraithConfig.BindPFlag("gitlab-api-token", scanGitlabCmd.Flags().Lookup("gitlab-api-token"))
	err = wraithConfig.BindPFlag("gitlab-projects", scanGitlabCmd.Flags().Lookup("gitlab-projects"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
