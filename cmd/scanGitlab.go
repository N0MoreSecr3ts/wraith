// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"github.com/spf13/viper"
	"time"
	"wraith/core"
	"wraith/version"

	"fmt"
	"github.com/spf13/cobra"
)

var viperScanGitlab *viper.Viper

// scanGitlabCmd represents the scanGitlab command
var scanGitlabCmd = &cobra.Command{
	Use:   "scanGitlab",
	Short: "Scan one or more gitlab groups or users for secrets",
	Long:  `Scan one or more gitlab groups or users for secrets`,
	Run: func(cmd *cobra.Command, args []string) {

		scanType := "gitlab"
		sess := core.NewSession(viperScanGitlab, scanType)

		//sess.Out.Info("%s\n\n", common.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", sess.BindAddress, sess.BindPort)

		// TODO need to validate this
		sess.GitlabAccessToken = viperScanGitlab.GetString("gitlab-api-token")

		sess.InitGitClient()

		core.GatherTargets(sess)
		core.GatherGitlabRepositories(sess)
		core.AnalyzeRepositories(sess)
		sess.Finish()

		if sess.JSONOutput || sess.CSVOutput {
			core.WriteOutput(sess)
		}

		core.PrintSessionStats(sess)

		if !sess.Silent {
			sess.Out.Important("%s", core.GitLabTanuki)
			sess.Out.Important("Press Ctrl+C to stop web server and exit.")
			select {}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanGitlabCmd)

	viperScanGitlab = core.SetConfig()

	scanGitlabCmd.Flags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	scanGitlabCmd.Flags().Int("bind-port", 9393, "The port for the webserver")
	scanGitlabCmd.Flags().Bool("csv", false, "Write results to --output-file in CSV format")
	scanGitlabCmd.Flags().Int("commit-depth", 0, "Set the depth for commits")
	scanGitlabCmd.Flags().Bool("debug", false, "Print debugging information")
	scanGitlabCmd.Flags().String("gitlab-api-token", "", "API token for access to Gitlab, see doc for necessary scope")
	scanGitlabCmd.Flags().String("gitlab-targets", "", "A space separated list of Gitlab users, projects or groups to scan")
	scanGitlabCmd.Flags().Bool("hide-secrets", false, "Hide secrets from output")
	scanGitlabCmd.Flags().String("ignore-extension", "", "a comma separated list of extensions to ignore")
	scanGitlabCmd.Flags().String("ignore-path", "", "a comma separated list of paths to ignore")
	scanGitlabCmd.Flags().Bool("in-mem-clone", false, "Clone repos in memory")
	scanGitlabCmd.Flags().Bool("json", false, "Write results to --output-file in JSON format")
	scanGitlabCmd.Flags().Int("match-level", 3, "Signature match level")
	scanGitlabCmd.Flags().Int("max-file-size", 50, "Max file size to scan")
	scanGitlabCmd.Flags().Bool("no-expand-orgs", false, "Don't add members to targets when processing organizations")
	scanGitlabCmd.Flags().Int("num-threads", 0, "The number of threads to execute with")
	scanGitlabCmd.Flags().String("output-dir", "$HOME/.wraith/", "Write csv and/or json files to directory")
	scanGitlabCmd.Flags().String("output-prefix", "wraith", "Prefix to prepend to datetime stamp for output files")
	scanGitlabCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanGitlabCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default.yml", "file(s) containing detection signatures.")
	scanGitlabCmd.Flags().Bool("silent", false, "No output")


	err := viperScanGitlab.BindPFlag("bind-address", scanGitlabCmd.Flags().Lookup("bind-address"))
	err = viperScanGitlab.BindPFlag("bind-port", scanGitlabCmd.Flags().Lookup("bind-port"))
	err = viperScanGitlab.BindPFlag("commit-depth", scanGitlabCmd.Flags().Lookup("commit-depth"))
	err = viperScanGitlab.BindPFlag("csv", scanGitlabCmd.Flags().Lookup("csv"))
	err = viperScanGitlab.BindPFlag("debug", scanGitlabCmd.Flags().Lookup("debug"))
	err = viperScanGitlab.BindPFlag("gitlab-api-token", scanGitlabCmd.Flags().Lookup("gitlab-api-token"))
	err = viperScanGitlab.BindPFlag("gitlab-targets", scanGitlabCmd.Flags().Lookup("gitlab-targets"))
	err = viperScanGitlab.BindPFlag("hide-secrets", scanGitlabCmd.Flags().Lookup("hide-secrets"))
	err = viperScanGitlab.BindPFlag("ignore-extension", scanGitlabCmd.Flags().Lookup("ignore-extension"))
	err = viperScanGitlab.BindPFlag("ignore-path", scanGitlabCmd.Flags().Lookup("ignore-extension"))
	err = viperScanGitlab.BindPFlag("in-mem-clone", scanGitlabCmd.Flags().Lookup("in-mem-clone"))
	err = viperScanGitlab.BindPFlag("json", scanGitlabCmd.Flags().Lookup("json"))
	err = viperScanGitlab.BindPFlag("match-level", scanGitlabCmd.Flags().Lookup("match-level"))
	err = viperScanGitlab.BindPFlag("max-file-size", scanGitlabCmd.Flags().Lookup("max-file-size"))
	err = viperScanGitlab.BindPFlag("expand-orgs", scanGitlabCmd.Flags().Lookup("expand-orgs"))
	err = viperScanGitlab.BindPFlag("num-threads", scanGitlabCmd.Flags().Lookup("num-threads"))
	err = viperScanGitlab.BindPFlag("output-dir", scanGithubCmd.Flags().Lookup("output-dir"))
	err = viperScanGitlab.BindPFlag("output-prefix", scanGithubCmd.Flags().Lookup("output-prefix"))
	err = viperScanGitlab.BindPFlag("scan-tests", scanGitlabCmd.Flags().Lookup("scan-tests"))
	err = viperScanGitlab.BindPFlag("signature-file", scanGitlabCmd.Flags().Lookup("signature-file"))
	err = viperScanGitlab.BindPFlag("silent", scanGitlabCmd.Flags().Lookup("silent"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
