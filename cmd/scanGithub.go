// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"time"
	"wraith/core"
	"wraith/version"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var viperScanGithub *viper.Viper

// scanGithubCmd represents the scanGithub command that will enumerate and scan github.com
var scanGithubCmd = &cobra.Command{
	Use:   "scanGithub",
	Short: "Scan one or more github.com orgs or users for secrets.",
	Long:  `Scan one or more github.com orgs or users for secrets.`,
	Run: func(cmd *cobra.Command, args []string) {

		scanType := "github"
		sess := core.NewSession(viperScanGithub, scanType)

		//sess.Out.Info("%s\n\n", common.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", sess.BindAddress, sess.BindPort)

		core.GatherTargets(sess)
		core.GatherRepositories(sess)
		core.AnalyzeRepositories(sess)
		sess.Finish()

		if sess.JSONOutput || sess.CSVOutput {
			core.WriteOutput(sess)
		}

		core.PrintSessionStats(sess)

		if !sess.Silent {
			sess.Out.Important("Press Ctrl+C to stop web server and exit.")
			select {}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanGithubCmd)

	viperScanGithub = core.SetConfig()

	scanGithubCmd.Flags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	scanGithubCmd.Flags().Int("bind-port", 9393, "The port for the webserver")
	scanGithubCmd.Flags().Int("commit-depth", 0, "Set the depth for commits")
	scanGithubCmd.Flags().Bool("csv", false, "Write results to --output-file in CSV format")
	scanGithubCmd.Flags().Bool("debug", false, "Print debugging information")
	scanGithubCmd.Flags().Bool("enterprise-scan", false, "Override github.com with a url specified via the --enterprise-url parameter")
	scanGithubCmd.Flags().String("enterprise-url", "", "Enterprise base url, must include protocol. IE https://github.org.com")
	scanGithubCmd.Flags().String("github-api-token", "", "API token for access to github, see doc for necessary scope")
	scanGithubCmd.Flags().String("github-targets", "", "A space separated list of github.com users or orgs to scan")
	scanGithubCmd.Flags().Bool("hide-secrets", false, "Hide secrets from output")
	scanGithubCmd.Flags().String("ignore-extension", "", "a comma separated list of extensions to ignore")
	scanGithubCmd.Flags().String("ignore-path", "", "a comma separated list of paths to ignore")
	scanGithubCmd.Flags().Bool("in-mem-clone", false, "Clone repos in memory")
	scanGithubCmd.Flags().Bool("json", false, "Write results to --output-file in JSON format")
	scanGithubCmd.Flags().Int("match-level", 3, "Signature match level")
	scanGithubCmd.Flags().Int("max-file-size", 50, "Max file size to scan")
	scanGithubCmd.Flags().Bool("no-expand-orgs", false, "Don't add members to targets when processing organizations")
	scanGithubCmd.Flags().Int("num-threads", 0, "The number of threads to execute with")
	scanGithubCmd.Flags().String("output-dir", "./", "Write csv and/or json files to directory")
	scanGithubCmd.Flags().String("output-prefix", "wraith", "Prefix to prepend to datetime stamp for output files")
	scanGithubCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanGithubCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default.yml", "file(s) containing detection signatures.")
	scanGithubCmd.Flags().Bool("silent", false, "No output")

	err := viperScanGithub.BindPFlag("bind-address", scanGithubCmd.Flags().Lookup("bind-address"))
	err = viperScanGithub.BindPFlag("bind-port", scanGithubCmd.Flags().Lookup("bind-port"))
	err = viperScanGithub.BindPFlag("commit-depth", scanGithubCmd.Flags().Lookup("commit-depth"))
	err = viperScanGithub.BindPFlag("csv", scanGithubCmd.Flags().Lookup("csv"))
	err = viperScanGithub.BindPFlag("debug", scanGithubCmd.Flags().Lookup("debug"))
	err = viperScanGithub.BindPFlag("enterprise-scan", scanGithubCmd.Flags().Lookup("enterprise-scan"))
	err = viperScanGithub.BindPFlag("enterprise-url", scanGithubCmd.Flags().Lookup("enterprise-url"))
	err = viperScanGithub.BindPFlag("github-api-token", scanGithubCmd.Flags().Lookup("github-api-token"))
	err = viperScanGithub.BindPFlag("github-targets", scanGithubCmd.Flags().Lookup("github-targets"))
	err = viperScanGithub.BindPFlag("hide-secrets", scanGithubCmd.Flags().Lookup("hide-secrets"))
	err = viperScanGithub.BindPFlag("ignore-extension", scanGithubCmd.Flags().Lookup("ignore-extension"))
	err = viperScanGithub.BindPFlag("ignore-path", scanGithubCmd.Flags().Lookup("ignore-extension"))
	err = viperScanGithub.BindPFlag("in-mem-clone", scanGithubCmd.Flags().Lookup("in-mem-clone"))
	err = viperScanGithub.BindPFlag("json", scanGithubCmd.Flags().Lookup("json"))
	err = viperScanGithub.BindPFlag("match-level", scanGithubCmd.Flags().Lookup("match-level"))
	err = viperScanGithub.BindPFlag("max-file-size", scanGithubCmd.Flags().Lookup("max-file-size"))
	err = viperScanGithub.BindPFlag("no-expand-orgs", scanGithubCmd.Flags().Lookup("no-expand-orgs"))
	err = viperScanGithub.BindPFlag("num-threads", scanGithubCmd.Flags().Lookup("num-threads"))
	err = viperScanGithub.BindPFlag("output-dir", scanGithubCmd.Flags().Lookup("output-dir"))
	err = viperScanGithub.BindPFlag("output-prefix", scanGithubCmd.Flags().Lookup("output-prefix"))
	err = viperScanGithub.BindPFlag("scan-tests", scanGithubCmd.Flags().Lookup("scan-tests"))
	err = viperScanGithub.BindPFlag("signature-file", scanGithubCmd.Flags().Lookup("signature-file"))
	err = viperScanGithub.BindPFlag("silent", scanGithubCmd.Flags().Lookup("silent"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
