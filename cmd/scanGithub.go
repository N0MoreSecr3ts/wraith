// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
	"wraith/core"
	"wraith/version"
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
		sess.GithubAccessToken = core.CheckGithubAPIToken(viperScanGithub.GetString("github-api-token"), sess)

		//sess.Out.Info("%s\n\n", common.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", "127.0.0.1", 9393)

		sess.Out.Debug("We have these orgs: %s\n", sess.UserOrgs)
		sess.Out.Debug("We have these users: %s\n", sess.UserLogins)
		sess.Out.Debug("We have these repos: %s\n", sess.UserRepos)

		//Create a github client to be used for the session
		sess.InitGitClient()

		// If we have github users and no orgs or repos then we default to scan
		// the visible repos of that user.
		if sess.UserLogins != nil && sess.UserOrgs == nil && sess.UserRepos == nil {
			//fmt.Println("I should not be here in the user") // TODO remove me
			//if sess.UserOrgs == nil && sess.UserRepos == nil {
			//sess.Out.Debug("I am scanning all the repos for the user(s) %s.\n", sess.UserLogins) // TODO remove me
			core.GatherUsers(sess)
			core.GetGithubRepositoriesFromOwner(sess)
			//}
		} else if sess.UserOrgs != nil && sess.UserLogins == nil && sess.UserRepos == nil {
			//fmt.Println("I should not be here in the org") // TODO remove me
			// If the user has only given orgs then we grab all te repos from those orgs
			//if sess.UserLogins == nil && sess.UserRepos == nil {
			//sess.Out.Debug("I am scanning all the repos for the org(s) %s.\n", sess.UserOrgs) // TODO remove me
			core.GatherOrgs(sess)
			core.GatherGithubOrgRepositories(sess)
			//}
		} else if sess.UserRepos != nil && sess.UserOrgs != nil {
			//fmt.Println("I should be here in the first step") // TODO remove me
			// If we have repo(s) given we need to ensure that we also have orgs or users. Wraith will then
			// look for the repo in the user or login lists and scan it.
			//if sess.UserOrgs != nil {
			//sess.Out.Debug("I am scanning the repo(s) %s in the org(s) %s.\n", sess.UserRepos, sess.UserOrgs) // TODO remove me
			core.GatherOrgs(sess)
			core.GatherGithubOrgRepositories(sess)
		} else if sess.UserRepos != nil && sess.UserLogins != nil {
			//sess.Out.Debug("I am scanning the repo(s) %s for the user(s) %s.\n", sess.UserRepos, sess.UserLogins) // TODO remove me
			core.GatherUsers(sess)
			core.GetGithubRepositoriesFromOwner(sess)
		} else {
			sess.Out.Error("You need to specify an org or user that contains the repo(s).\n")
		}

		core.AnalyzeRepositories(sess)
		sess.Finish()

		core.PrintSessionStats(sess)

		if !sess.Silent {
			sess.Out.Important("Press Ctrl+C to stop web server and exit.\n")
			select {}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanGithubCmd)

	viperScanGithub = core.SetConfig()

	scanGithubCmd.Flags().Bool("expand-orgs", false, "Add members to targets when processing organizations")
	scanGithubCmd.Flags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	scanGithubCmd.Flags().Int("bind-port", 9393, "The port for the webserver")
	scanGithubCmd.Flags().Int("commit-depth", 0, "Set the depth for commits")
	scanGithubCmd.Flags().Bool("debug", false, "Print debugging information")
	scanGithubCmd.Flags().String("github-api-token", "", "API token for access to github, see doc for necessary scope")
	scanGithubCmd.Flags().String("github-targets", "", "A space separated list of github.com users or orgs to scan")
	scanGithubCmd.Flags().Bool("hide-secrets", false, "Hide secrets from output")
	scanGithubCmd.Flags().String("ignore-extension", "", "a comma separated list of extensions to ignore")
	scanGithubCmd.Flags().String("ignore-path", "", "a comma separated list of paths to ignore")
	scanGithubCmd.Flags().Bool("in-mem-clone", false, "Clone repos in memory")
	scanGithubCmd.Flags().Int("match-level", 3, "Signature match level")
	scanGithubCmd.Flags().Int("max-file-size", 50, "Max file size to scan")
	scanGithubCmd.Flags().Int("num-threads", 0, "The number of threads to execute with")
	scanGithubCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanGithubCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default.yml", "file(s) containing detection signatures.")
	scanGithubCmd.Flags().Bool("silent", false, "No output")
	scanGithubCmd.Flags().String("github-url", "", "The api endpoint for github.com")
	scanGithubCmd.Flags().String("github-orgs", "", "A coma separated list of github orgs to scan")
	scanGithubCmd.Flags().String("github-repos", "", "A coma separated list of github repositories to scan")
	scanGithubCmd.Flags().String("github-users", "", "A coma separated list of github.com users to scan")

	err := viperScanGithub.BindPFlag("bind-address", scanGithubCmd.Flags().Lookup("bind-address"))
	err = viperScanGithub.BindPFlag("github-url", scanGithubCmd.Flags().Lookup("github-url"))
	err = viperScanGithub.BindPFlag("bind-port", scanGithubCmd.Flags().Lookup("bind-port"))
	err = viperScanGithub.BindPFlag("commit-depth", scanGithubCmd.Flags().Lookup("commit-depth"))
	err = viperScanGithub.BindPFlag("debug", scanGithubCmd.Flags().Lookup("debug"))
	err = viperScanGithub.BindPFlag("enterprise-scan", scanGithubCmd.Flags().Lookup("enterprise-scan"))
	err = viperScanGithub.BindPFlag("github-api-token", scanGithubCmd.Flags().Lookup("github-api-token"))
	err = viperScanGithub.BindPFlag("github-targets", scanGithubCmd.Flags().Lookup("github-targets"))
	err = viperScanGithub.BindPFlag("hide-secrets", scanGithubCmd.Flags().Lookup("hide-secrets"))
	err = viperScanGithub.BindPFlag("ignore-extension", scanGithubCmd.Flags().Lookup("ignore-extension"))
	err = viperScanGithub.BindPFlag("ignore-path", scanGithubCmd.Flags().Lookup("ignore-extension"))
	err = viperScanGithub.BindPFlag("in-mem-clone", scanGithubCmd.Flags().Lookup("in-mem-clone"))
	err = viperScanGithub.BindPFlag("match-level", scanGithubCmd.Flags().Lookup("match-level"))
	err = viperScanGithub.BindPFlag("max-file-size", scanGithubCmd.Flags().Lookup("max-file-size"))
	err = viperScanGithub.BindPFlag("expand-orgs", scanGithubCmd.Flags().Lookup("expand-orgs"))
	err = viperScanGithub.BindPFlag("num-threads", scanGithubCmd.Flags().Lookup("num-threads"))
	err = viperScanGithub.BindPFlag("scan-tests", scanGithubCmd.Flags().Lookup("scan-tests"))
	err = viperScanGithub.BindPFlag("signature-file", scanGithubCmd.Flags().Lookup("signature-file"))
	err = viperScanGithub.BindPFlag("silent", scanGithubCmd.Flags().Lookup("silent"))
	err = viperScanGithub.BindPFlag("github-orgs", scanGithubCmd.Flags().Lookup("github-orgs"))
	err = viperScanGithub.BindPFlag("github-repos", scanGithubCmd.Flags().Lookup("github-repos"))
	err = viperScanGithub.BindPFlag("github-users", scanGithubCmd.Flags().Lookup("github-users"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
