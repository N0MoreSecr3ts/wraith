// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/N0MoreSecr3ts/wraith/core"
	"github.com/N0MoreSecr3ts/wraith/version"

	"github.com/spf13/cobra"
)

// scanGithubEnterpriseCmd represents the scanGithubEnterprise command
var scanGithubEnterpriseCmd = &cobra.Command{
	Use:   "scanGithubEnterprise",
	Short: "Scan one or more github enterprise organizations and repos for secrets.",
	Long:  "Scan one or more github enterprise organizations and repos for secrets. - v" + version.AppVersion(),
	Run: func(cmd *cobra.Command, args []string) {

		// Set the scan type and start a new session
		scanType := "github-enterprise"
		sess := core.NewSession(wraithConfig, scanType)

		// Ensure user input exists and validate it
		sess.ValidateUserInput(wraithConfig)

		// Check for a token. If no token is present we should default to scan but give a message
		// that no token is available so only public repos will be scanned
		sess.GithubAccessToken = core.CheckGithubAPIToken(wraithConfig.GetString("github-enterprise-api-token"), sess)

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

		//Create a github client to be used for the session
		sess.InitGitClient()

		// If we have github users and no orgs or repos then we default to scan
		// the visible repos of that user.
		if sess.UserLogins != nil {
			if sess.UserOrgs == nil && sess.UserRepos == nil {
				core.GatherUsers(sess)
			}
		}

		// If the user has only given orgs then we grab all te repos from those orgs
		if sess.UserOrgs != nil {
			if sess.UserLogins == nil && sess.UserRepos == nil {
				core.GatherOrgs(sess)
			}
		}

		// If we have repo(s) given we need to ensure that we also have orgs or users. Wraith will then
		// look for the repo in the user or login lists and scan it.
		if sess.UserRepos != nil {
			if sess.UserOrgs != nil {
				core.GatherOrgs(sess)
				core.GatherGithubOrgRepositories(sess)
			} else if sess.UserLogins != nil {
				core.GatherUsers(sess)
				core.GatherGithubRepositoriesFromOwner(sess)
			} else {
				sess.Out.Error("You need to specify an org or user that contains the repo(s).\n")
				os.Exit(1)
			}
		}

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
	rootCmd.AddCommand(scanGithubEnterpriseCmd)

	scanGithubEnterpriseCmd.Flags().Bool("add-org-members", false, "Add members to targets when processing organizations")
	scanGithubEnterpriseCmd.Flags().String("github-enterprise-api-token", "", "API token for github access, see documentation for necessary scope")
	scanGithubEnterpriseCmd.Flags().StringSlice("github-enterprise-orgs", nil, "List of github orgs to scan")
	scanGithubEnterpriseCmd.Flags().StringSlice("github-enterprise-repos", nil, "List of github repositories to scan")
	scanGithubEnterpriseCmd.Flags().String("github-enterprise-url", "", "Entperise Github instance. I.E. https://github.org.com")
	scanGithubEnterpriseCmd.Flags().StringSlice("github-enterprise-users", nil, "List of github.com users to scan")
	scanGithubEnterpriseCmd.Flags().Float64("commit-depth", -1, "Set the commit depth to scan")

	err := wraithConfig.BindPFlag("add-org-members", scanGithubEnterpriseCmd.Flags().Lookup("add-org-members"))
	err = wraithConfig.BindPFlag("commit-depth", scanGithubEnterpriseCmd.Flags().Lookup("commit-depth"))
	err = wraithConfig.BindPFlag("github-enterprise-api-token", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-api-token"))
	err = wraithConfig.BindPFlag("github-enterprise-orgs", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-orgs"))
	err = wraithConfig.BindPFlag("github-enterprise-repos", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-repos"))
	err = wraithConfig.BindPFlag("github-enterprise-url", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-url"))
	err = wraithConfig.BindPFlag("github-enterprise-users", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-users"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
