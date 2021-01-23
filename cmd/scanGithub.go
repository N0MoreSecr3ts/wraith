// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
	"wraith/core"
)

var viperScanGithub *viper.Viper

// scanGithubCmd represents the scanGithub command that will enumerate and scan github.com
var scanGithubCmd = &cobra.Command{
	Use:   "scanGithub",
	Short: "Scan one or more github.com orgs or users for secrets.",
	Long:  `Scan one or more github.com orgs or users for secrets.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Set the scan type and start a new session
		scanType := "github"
		sess := core.NewSession(viperScanGithub, scanType)

		// Ensure user input exists and validate it
		sess.ValidateUserInput(viperScanGithub)

		// Check for a token. If no token is present we should default to scan but give a message
		// that no token is available so only public repos will be scanned
		// TODO do not exit out if no token but drop a message saying only public repos will be scanned
		sess.GithubAccessToken = core.CheckGithubAPIToken(viperScanGithub.GetString("github-api-token"), sess)

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
		if sess.UserLogins != nil && sess.UserOrgs == nil && sess.UserRepos == nil {

			core.GatherUsers(sess)
			core.GatherGithubRepositoriesFromOwner(sess)

		} else if sess.UserOrgs != nil && sess.UserLogins == nil && sess.UserRepos == nil {

			// If the user has only given orgs then we grab all te repos from those orgs
			core.GatherOrgs(sess)
			core.GatherGithubOrgRepositories(sess)

		} else if sess.UserRepos != nil && sess.UserOrgs != nil {

			// If we have repo(s) given we need to ensure that we also have orgs or users. Wraith will then
			// look for the repo in the user or login lists and scan it.
			core.GatherOrgs(sess)
			core.GatherGithubOrgRepositories(sess)
		} else if sess.UserRepos != nil && sess.UserLogins != nil {

			// If we have repos and users then we are going to scan for the repo in that users account
			core.GatherUsers(sess)
			core.GatherGithubRepositoriesFromOwner(sess)
		} else {
			// Catchall for not being able to scan any as either we have no information or
			// we don't have the rights kinds of information
			sess.Out.Error("You need to specify an org or user that contains the repo(s).\n")
			os.Exit(1)
		}

		core.AnalyzeRepositories(sess)
		sess.Finish()

		core.PrintSessionStats(sess)

		if !sess.Silent && sess.WebServer {
			sess.Out.Important("Press Ctrl+C to stop web server and exit.\n")
			select {}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanGithubCmd)

	viperScanGithub = core.SetConfig()

	scanGithubCmd.Flags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	scanGithubCmd.Flags().Int("bind-port", 9393, "The port for the webserver")
	scanGithubCmd.Flags().Int("confidence-level", 3, "The confidence level level of the expressions used to find matches")
	scanGithubCmd.Flags().Float64("commit-depth", -1, "Set the commit depth to scan")
	scanGithubCmd.Flags().Bool("debug", false, "Print available debugging information to stdout")
	scanGithubCmd.Flags().Bool("add-org-members", false, "Add members to targets when processing organizations")
	scanGithubCmd.Flags().String("github-api-token", "", "API token for github access, see documentation for necessary scope")
	scanGithubCmd.Flags().StringSlice("github-orgs", nil, "List of github orgs to scan")
	scanGithubCmd.Flags().StringSlice("github-repos", nil, "List of github repositories to scan")
	scanGithubCmd.Flags().StringSlice("github-users", nil, "List of github.com users to scan")
	scanGithubCmd.Flags().Bool("hide-secrets", false, "Do not print secrets to any supported output")
	scanGithubCmd.Flags().StringSlice("ignore-extension", nil, "List of file extensions to ignore")
	scanGithubCmd.Flags().StringSlice("ignore-path", nil, "List of file paths to ignore")
	scanGithubCmd.Flags().Int("max-file-size", 10, "Max file size to scan (in MB)")
	scanGithubCmd.Flags().Int("num-threads", -1, "Number of execution threads")
	scanGithubCmd.Flags().Bool("scan-forks", false, "Scan repositories forked by users or orgs")
	scanGithubCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanGithubCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default.yaml", "file(s) containing detection signatures.")
	scanGithubCmd.Flags().String("signature-path", "$HOME/.wraith/signatures", "path containing detection signatures.")
	scanGithubCmd.Flags().Bool("silent", false, "Suppress all output except for errors")
	scanGithubCmd.Flags().Bool("web-server", false, "Enable the web interface for scan output")

	err := viperScanGithub.BindPFlag("bind-address", scanGithubCmd.Flags().Lookup("bind-address"))
	err = viperScanGithub.BindPFlag("bind-port", scanGithubCmd.Flags().Lookup("bind-port"))
	err = viperScanGithub.BindPFlag("commit-depth", scanGithubCmd.Flags().Lookup("commit-depth"))
	err = viperScanGithub.BindPFlag("debug", scanGithubCmd.Flags().Lookup("debug"))
	err = viperScanGithub.BindPFlag("add-org-members", scanGithubCmd.Flags().Lookup("add-org-members"))
	err = viperScanGithub.BindPFlag("github-api-token", scanGithubCmd.Flags().Lookup("github-api-token"))
	err = viperScanGithub.BindPFlag("hide-secrets", scanGithubCmd.Flags().Lookup("hide-secrets"))
	err = viperScanGithub.BindPFlag("ignore-extension", scanGithubCmd.Flags().Lookup("ignore-extension"))
	err = viperScanGithub.BindPFlag("ignore-path", scanGithubCmd.Flags().Lookup("ignore-path"))
	err = viperScanGithub.BindPFlag("confidence-level", scanGithubCmd.Flags().Lookup("confidence-level"))
	err = viperScanGithub.BindPFlag("max-file-size", scanGithubCmd.Flags().Lookup("max-file-size"))
	err = viperScanGithub.BindPFlag("num-threads", scanGithubCmd.Flags().Lookup("num-threads"))
	err = viperScanGithub.BindPFlag("scan-forks", scanGithubCmd.Flags().Lookup("scan-forks"))
	err = viperScanGithub.BindPFlag("scan-tests", scanGithubCmd.Flags().Lookup("scan-tests"))
	err = viperScanGithub.BindPFlag("signature-file", scanGithubCmd.Flags().Lookup("signature-file"))
	err = viperScanGithub.BindPFlag("signature-path", scanGithubCmd.Flags().Lookup("signature-path"))
	err = viperScanGithub.BindPFlag("silent", scanGithubCmd.Flags().Lookup("silent"))
	err = viperScanGithub.BindPFlag("github-orgs", scanGithubCmd.Flags().Lookup("github-orgs"))
	err = viperScanGithub.BindPFlag("github-repos", scanGithubCmd.Flags().Lookup("github-repos"))
	err = viperScanGithub.BindPFlag("github-users", scanGithubCmd.Flags().Lookup("github-users"))
	err = viperScanGithub.BindPFlag("web-server", scanGithubCmd.Flags().Lookup("web-server"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
