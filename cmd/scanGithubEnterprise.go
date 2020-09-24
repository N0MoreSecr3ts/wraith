// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"time"
	"wraith/core"
	"wraith/version"

	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

// viperScanGithubEnterprise holds the configuration data for this subcommand
var viperScanGithubEnterprise *viper.Viper

// scanGithubEnterpriseCmd represents the scanGithubEnterprise command
var scanGithubEnterpriseCmd = &cobra.Command{
	Use:   "scanGithubEnterprise",
	Short: "Scan one or more github enterprise organizations and repos for secrets.",
	Long:  "Scan one or more github enterprise organizations and repos for secrets. - v" + version.AppVersion(),
	Run: func(cmd *cobra.Command, args []string) {

		// Set the scan type and start a new session
		scanType := "github-enterprise"
		sess := core.NewSession(viperScanGithubEnterprise, scanType)

		// Ensure user input exists and validate it
		sess.ValidateUserInput(viperScanGithubEnterprise)

		// Check for a token. If no token is present we should default to scan but give a message
		// that no token is available so only public repos will be scanned
		sess.GithubAccessToken = core.CheckGithubAPIToken(viperScanGithubEnterprise.GetString("github-api-token"), sess)

		//sess.Out.Info("%s\n\n", common.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", sess.BindAddress, sess.BindPort)

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

		core.PrintSessionStats(sess)

		if !sess.Silent {
			sess.Out.Important("Press Ctrl+C to stop web server and exit.\n")
			select {}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanGithubEnterpriseCmd)

	viperScanGithubEnterprise = core.SetConfig()

	scanGithubEnterpriseCmd.Flags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	scanGithubEnterpriseCmd.Flags().Int("bind-port", 9393, "The port for the webserver")
	scanGithubEnterpriseCmd.Flags().Int("confidence-level", 3, "The confidence level level of the expressions used to find matches")
	scanGithubEnterpriseCmd.Flags().Float64("commit-depth", -1, "Set the commit depth to scan")
	scanGithubEnterpriseCmd.Flags().Bool("debug", false, "Print debugging information")
	scanGithubEnterpriseCmd.Flags().Bool("gather-org-members", false, "Add members to targets when processing organizations")
	scanGithubEnterpriseCmd.Flags().String("github-enterprise-api-token", "", "API token for access to github, see doc for necessary scope")
	scanGithubEnterpriseCmd.Flags().StringSlice("github-enterprise-orgs", nil, "List of github orgs to scan")
	scanGithubEnterpriseCmd.Flags().StringSlice("github-enterprise-repos", nil, "List of github repositories to scan")
	scanGithubEnterpriseCmd.Flags().StringSlice("github-enterprise-users", nil, "List of github.com users to scan")
	scanGithubEnterpriseCmd.Flags().Bool("hide-secrets", false, "Hide secrets from any supported output")
	scanGithubEnterpriseCmd.Flags().StringSlice("ignore-extension", nil, "List of extensions to ignore")
	scanGithubEnterpriseCmd.Flags().StringSlice("ignore-path", nil, "List of paths to ignore")
	//scanGithubEnterpriseCmd.Flags().Bool("in-mem-clone", false, "Clone repos in memory")
	scanGithubEnterpriseCmd.Flags().Int("max-file-size", 10, "Max file size to scan in MB")
	scanGithubEnterpriseCmd.Flags().Int("num-threads", -1, "Number of threads to execute with")
	scanGithubEnterpriseCmd.Flags().Bool("scan-forks", true, "Scan repositories forked by users or orgs")
	scanGithubEnterpriseCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanGithubEnterpriseCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default.yml", "file(s) containing detection signatures.")
	scanGithubEnterpriseCmd.Flags().String("signature-path", "$HOME/.wraith/signatures", "path containing detection signatures.")
	scanGithubEnterpriseCmd.Flags().Bool("silent", false, "Suppress all output except for errors")

	err := viperScanGithubEnterprise.BindPFlag("bind-address", scanGithubEnterpriseCmd.Flags().Lookup("bind-address"))
	err = viperScanGithubEnterprise.BindPFlag("bind-port", scanGithubEnterpriseCmd.Flags().Lookup("bind-port"))
	err = viperScanGithubEnterprise.BindPFlag("commit-depth", scanGithubEnterpriseCmd.Flags().Lookup("commit-depth"))
	err = viperScanGithubEnterprise.BindPFlag("debug", scanGithubEnterpriseCmd.Flags().Lookup("debug"))
	err = viperScanGithubEnterprise.BindPFlag("gather-org-members", scanGithubEnterpriseCmd.Flags().Lookup("gather-org-members"))
	err = viperScanGithubEnterprise.BindPFlag("github-enterprise-api-token", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-api-token"))
	err = viperScanGithubEnterprise.BindPFlag("hide-secrets", scanGithubEnterpriseCmd.Flags().Lookup("hide-secrets"))
	err = viperScanGithubEnterprise.BindPFlag("ignore-extension", scanGithubEnterpriseCmd.Flags().Lookup("ignore-extension"))
	err = viperScanGithubEnterprise.BindPFlag("ignore-path", scanGithubEnterpriseCmd.Flags().Lookup("ignore-extension"))
	//err = viperScanGithubEnterprise.BindPFlag("in-mem-clone", scanGithubEnterpriseCmd.Flags().Lookup("in-mem-clone"))
	err = viperScanGithubEnterprise.BindPFlag("confidence-level", scanGithubEnterpriseCmd.Flags().Lookup("confidence-level"))
	err = viperScanGithubEnterprise.BindPFlag("max-file-size", scanGithubEnterpriseCmd.Flags().Lookup("max-file-size"))
	err = viperScanGithubEnterprise.BindPFlag("num-threads", scanGithubEnterpriseCmd.Flags().Lookup("num-threads"))
	err = viperScanGithubEnterprise.BindPFlag("scan-forks", scanGithubEnterpriseCmd.Flags().Lookup("scan-forks"))
	err = viperScanGithubEnterprise.BindPFlag("scan-tests", scanGithubEnterpriseCmd.Flags().Lookup("scan-tests"))
	err = viperScanGithubEnterprise.BindPFlag("signature-file", scanGithubEnterpriseCmd.Flags().Lookup("signature-file"))
	err = viperScanGithubEnterprise.BindPFlag("signature-path", scanGithubEnterpriseCmd.Flags().Lookup("signature-path"))
	err = viperScanGithubEnterprise.BindPFlag("silent", scanGithubEnterpriseCmd.Flags().Lookup("silent"))
	err = viperScanGithubEnterprise.BindPFlag("github-enterprise-orgs", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-orgs"))
	err = viperScanGithubEnterprise.BindPFlag("github-enterprise-repos", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-repos"))
	err = viperScanGithubEnterprise.BindPFlag("github-enterprise-users", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-users"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
