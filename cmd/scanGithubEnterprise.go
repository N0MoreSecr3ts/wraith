// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"os"
	"time"
	"wraith/core"
	"wraith/version"

	"github.com/spf13/viper"

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

		scanType := "github-enterprise"
		sess := core.NewSession(viperScanGithubEnterprise, scanType)

		sess.UserDirtyRepos = viperScanGithubEnterprise.GetString("github-enterprise-repos")
		sess.UserDirtyOrgs = viperScanGithubEnterprise.GetString("github-enterprise-orgs")
		sess.GithubAccessToken = core.CheckGithubAPIToken(viperScanGithubEnterprise.GetString("github-enterprise-api-token"), sess) //TODO can we clean this function up at all

		//fmt.Println( viperScanGithubEnterprise.GetString("github-enterprise-repos")) //TODO remove me
		//fmt.Println( viperScanGithubEnterprise.GetString("github-enterprise-orgs")) //TODO remove me

		if sess.UserDirtyRepos == "" && sess.UserDirtyOrgs == "" {
			fmt.Println("You must enter either an org or repo[s] to scan")
			os.Exit(2)
		}

		if sess.UserDirtyOrgs != "" {
			core.ValidateGHInput(sess)
		}

		fmt.Println(sess.UserOrgs)
		if len(sess.UserRepos) >= 1 && len(sess.UserOrgs) < 1 {
			fmt.Println("You need to specify an org that contains the repo(s).")
			os.Exit(2)
		}

		//sess.Out.Info("%s\n\n", common.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", sess.BindAddress, sess.BindPort)

		sess.GithubAccessToken = core.CheckGithubAPIToken(viperScanGithubEnterprise.GetString("github-enterprise-api-token"), sess)
		sess.InitGitClient()

		core.GatherOrgs(sess)
		//fmt.Println(sess.Organizations) //TODO remove me
		core.GatherGithubRepositories(sess)
		core.AnalyzeRepositories(sess)
		sess.Finish()

		core.PrintSessionStats(sess)

		if !sess.Silent {
			sess.Out.Important("Press Ctrl+C to stop web server and exit.")
			select {}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanGithubEnterpriseCmd)

	viperScanGithubEnterprise = core.SetConfig()

	scanGithubEnterpriseCmd.Flags().Bool("expand-orgs", false, "Add members to targets when processing organizations")
	scanGithubEnterpriseCmd.Flags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	scanGithubEnterpriseCmd.Flags().Int("bind-port", 9393, "The port for the webserver")
	scanGithubEnterpriseCmd.Flags().Bool("debug", false, "Print debugging information")
	scanGithubEnterpriseCmd.Flags().Bool("hide-secrets", false, "Hide secrets in any supported output")
	scanGithubEnterpriseCmd.Flags().Bool("json", false, "output json format")
	scanGithubEnterpriseCmd.Flags().Bool("load-triage", false, "load a triage file")
	scanGithubEnterpriseCmd.Flags().Bool("scan-forks", true, "Scan forked repositories")
	scanGithubEnterpriseCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanGithubEnterpriseCmd.Flags().Bool("silent", false, "Suppress all output except for errors")
	scanGithubEnterpriseCmd.Flags().Int64("max-file-size", 50, "Max file size to scan")
	scanGithubEnterpriseCmd.Flags().Int("commit-depth", 0, "The commit depth you want to travel to, 0=all")
	scanGithubEnterpriseCmd.Flags().Int("match-level", 3, "The match level level of the expressions used to find matches")
	scanGithubEnterpriseCmd.Flags().String("github-enterprise-api-token", "", "API token for access to github enterprise, see doc for necessary scope")
	//scanGithubEnterpriseCmd.Flags().String("github-targets", "", "A space separated list of github users or orgs to scan")
	scanGithubEnterpriseCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default.yml", "file(s) containing detection signatures.")
	scanGithubEnterpriseCmd.Flags().Int("num-threads", 0, "The number of threads to execute with")
	scanGithubEnterpriseCmd.Flags().Bool("in-mem-clone", false, "Clone repos in memory")
	scanGithubEnterpriseCmd.Flags().String("github-enterprise-url", "", "The api endpoint for github enterprise")
	scanGithubEnterpriseCmd.Flags().String("ignore-extension", "", "a list of extensions to ignore during a scan")
	scanGithubEnterpriseCmd.Flags().String("ignore-path", "", "a list of paths to ignore during a scan")
	scanGithubEnterpriseCmd.Flags().String("github-enterprise-orgs", "", "A coma separated list of github enterprise orgs to scan")
	scanGithubEnterpriseCmd.Flags().String("github-enterprise-repos", "", "A coma separated list of github enterprise repositories to scan")

	err := viperScanGithubEnterprise.BindPFlag("debug", scanGithubEnterpriseCmd.Flags().Lookup("debug"))
	err = viperScanGithubEnterprise.BindPFlag("hide-secrets", scanGithubEnterpriseCmd.Flags().Lookup("hide-secrets"))
	err = viperScanGithubEnterprise.BindPFlag("scan-tests", scanGithubEnterpriseCmd.Flags().Lookup("scan-tests"))
	err = viperScanGithubEnterprise.BindPFlag("silent", scanGithubEnterpriseCmd.Flags().Lookup("silent"))
	err = viperScanGithubEnterprise.BindPFlag("max-file-size", scanGithubEnterpriseCmd.Flags().Lookup("max-file-size"))
	err = viperScanGithubEnterprise.BindPFlag("commit-depth", scanGithubEnterpriseCmd.Flags().Lookup("commit-depth"))
	err = viperScanGithubEnterprise.BindPFlag("match-level", scanGithubEnterpriseCmd.Flags().Lookup("match-level"))
	err = viperScanGithubEnterprise.BindPFlag("github-enterprise-api-token", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-api-token"))
	err = viperScanGithubEnterprise.BindPFlag("github-enterprise-orgs", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-orgs"))
	err = viperScanGithubEnterprise.BindPFlag("github-enterprise-repos", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-repos"))
	err = viperScanGithubEnterprise.BindPFlag("github-enterprise-url", scanGithubEnterpriseCmd.Flags().Lookup("github-enterprise-url"))
	err = viperScanGithubEnterprise.BindPFlag("ignore-extension", scanGithubEnterpriseCmd.Flags().Lookup("ignore-extension"))
	err = viperScanGithubEnterprise.BindPFlag("ignore-path", scanGithubEnterpriseCmd.Flags().Lookup("ignore-path"))
	err = viperScanGithubEnterprise.BindPFlag("signature-file", scanGithubEnterpriseCmd.Flags().Lookup("signature-file"))
	err = viperScanGithubEnterprise.BindPFlag("bind-address", scanGithubCmd.Flags().Lookup("bind-address"))
	err = viperScanGithubEnterprise.BindPFlag("bind-port", scanGithubCmd.Flags().Lookup("bind-port"))
	//err = viperScanGithubEnterprise.BindPFlag("github-targets", scanGithubCmd.Flags().Lookup("github-targets"))
	err = viperScanGithubEnterprise.BindPFlag("in-mem-clone", scanGithubCmd.Flags().Lookup("in-mem-clone"))
	err = viperScanGithubEnterprise.BindPFlag("expand-orgs", scanGithubCmd.Flags().Lookup("expand-orgs"))
	err = viperScanGithubEnterprise.BindPFlag("num-threads", scanGithubCmd.Flags().Lookup("num-threads"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
