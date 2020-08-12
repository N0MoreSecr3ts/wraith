// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"github.com/spf13/viper"
	"os"
	"time"
	"wraith/core"
	"wraith/version"

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
		sess, err := core.NewSession(viperScanGitlab, scanType)

		if err != nil {
			sess.Out.Error("Failed to generate a new session: %s\n", err.Error())
			os.Exit(1)
		}

		//sess.Out.Info("%s\n\n", common.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", sess.BindAddress, sess.BindPort)

		core.GatherTargets(sess)
		core.GatherRepositories(sess)
		core.AnalyzeRepositories(sess)
		sess.Finish()

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

	scanGitlabCmd.Flags().Bool("debug", false, "Print debugging information")
	scanGitlabCmd.Flags().Bool("hide-secrets", false, "Hide secrets from output")
	scanGitlabCmd.Flags().Bool("in-mem-clone", false, "Clone repos in memory")
	scanGitlabCmd.Flags().Bool("no-expand-orgs", false, "Don't add members to targets when processing organizations")
	scanGitlabCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanGitlabCmd.Flags().Bool("silent", false, "No output")
	scanGitlabCmd.Flags().Int("bind-port", 9393, "The port for the webserver")
	scanGitlabCmd.Flags().Int("commit-depth", 0, "Set the depth for commits")
	scanGitlabCmd.Flags().Int("match-level", 3, "Signature match level")
	scanGitlabCmd.Flags().Int("max-file-size", 50, "Max file size to scan")
	scanGitlabCmd.Flags().Int("num-threads", 0, "The number of threads to execute with")
	scanGitlabCmd.Flags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	scanGitlabCmd.Flags().String("gitlab-api-token", "", "API token for access to Gitlab, see doc for necessary scope")
	scanGitlabCmd.Flags().String("gitlab-targets", "", "A space separated list of Gitlab users, projects or groups to scan")
	scanGitlabCmd.Flags().String("ignore-extension", "", "a comma separated list of extensions to ignore")
	scanGitlabCmd.Flags().String("ignore-path", "", "a comma separated list of paths to ignore")
	scanGitlabCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default_signatures.yml", "file(s) containing detection signatures.")

	viperScanGitlab.BindPFlag("bind-address", scanGitlabCmd.Flags().Lookup("bind-address"))
	viperScanGitlab.BindPFlag("bind-port", scanGitlabCmd.Flags().Lookup("bind-port"))
	viperScanGitlab.BindPFlag("commit-depth", scanGitlabCmd.Flags().Lookup("commit-depth"))
	viperScanGitlab.BindPFlag("debug", scanGitlabCmd.Flags().Lookup("debug"))
	viperScanGitlab.BindPFlag("gitlab-api-token", scanGitlabCmd.Flags().Lookup("gitlab-api-token"))
	viperScanGitlab.BindPFlag("gitlab-targets", scanGitlabCmd.Flags().Lookup("gitlab-targets"))
	viperScanGitlab.BindPFlag("hide-secrets", scanGitlabCmd.Flags().Lookup("hide-secrets"))
	viperScanGitlab.BindPFlag("ignore-extension", scanGitlabCmd.Flags().Lookup("ignore-extension"))
	viperScanGitlab.BindPFlag("ignore-path", scanGitlabCmd.Flags().Lookup("ignore-extension"))
	viperScanGitlab.BindPFlag("in-mem-clone", scanGitlabCmd.Flags().Lookup("in-mem-clone"))
	viperScanGitlab.BindPFlag("match-level", scanGitlabCmd.Flags().Lookup("match-level"))
	viperScanGitlab.BindPFlag("max-file-size", scanGitlabCmd.Flags().Lookup("max-file-size"))
	viperScanGitlab.BindPFlag("no-expand-orgs", scanGitlabCmd.Flags().Lookup("no-expand-orgs"))
	viperScanGitlab.BindPFlag("num-threads", scanGitlabCmd.Flags().Lookup("num-threads"))
	viperScanGitlab.BindPFlag("scan-tests", scanGitlabCmd.Flags().Lookup("scan-tests"))
	viperScanGitlab.BindPFlag("signature-file", scanGitlabCmd.Flags().Lookup("signature-file"))
	viperScanGitlab.BindPFlag("silent", scanGitlabCmd.Flags().Lookup("silent"))

}
