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

var viperScanLocalGitRepo *viper.Viper

// scanLocalGitRepoCmd represents the scanLocalGitRepo command
var scanLocalGitRepoCmd = &cobra.Command{
	Use:   "scanLocalGitRepo",
	Short: "Scan a git repo on a local machine",
	Long:  "Scan a git repo on a local machine",
	Run: func(cmd *cobra.Command, args []string) {

		scanType := "localGit"
		sess := core.NewSession(viperScanLocalGitRepo, scanType)

		//sess.Out.Info("%s\n\n", common.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", "127.0.0.1", 9393)

		core.GatherLocalRepositories(sess)
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
	rootCmd.AddCommand(scanLocalGitRepoCmd)

	viperScanLocalGitRepo = core.SetConfig()

	scanLocalGitRepoCmd.Flags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	scanLocalGitRepoCmd.Flags().Int("bind-port", 9393, "The port for the webserver")
	scanLocalGitRepoCmd.Flags().Int("commit-depth", 0, "Set the depth for commits")
	scanLocalGitRepoCmd.Flags().Bool("csv", false, "Write results to --output-file in CSV format")
	scanLocalGitRepoCmd.Flags().Bool("debug", false, "Print debugging information")
	scanLocalGitRepoCmd.Flags().Bool("hide-secrets", false, "Hide secrets from output")
	scanLocalGitRepoCmd.Flags().String("ignore-extension", "", "a comma separated list of extensions to ignore")
	scanLocalGitRepoCmd.Flags().String("ignore-path", "", "a comma separated list of paths to ignore")
	scanLocalGitRepoCmd.Flags().Bool("in-mem-clone", false, "Clone repos in memory")
	scanLocalGitRepoCmd.Flags().Bool("json", false, "Write results to --output-file in JSON format")
	scanLocalGitRepoCmd.Flags().String("local-dirs", "", "local disk parent dir containing git repos")
	scanLocalGitRepoCmd.Flags().Int("match-level", 3, "Signature match level")
	scanLocalGitRepoCmd.Flags().Int("max-file-size", 50, "Max file size to scan")
	scanLocalGitRepoCmd.Flags().Bool("no-expand-orgs", false, "Don't add members to targets when processing organizations")
	scanLocalGitRepoCmd.Flags().Int("num-threads", 0, "The number of threads to execute with")
	scanLocalGitRepoCmd.Flags().String("output-dir", "./", "Write csv and/or json files to directory")
	scanLocalGitRepoCmd.Flags().String("output-prefix", "wraith", "Prefix to prepend to datetime stamp for output files")
	scanLocalGitRepoCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanLocalGitRepoCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default.yml", "file(s) containing detection signatures.")
	scanLocalGitRepoCmd.Flags().Bool("silent", false, "No output")


	err := viperScanLocalGitRepo.BindPFlag("bind-address", scanLocalGitRepoCmd.Flags().Lookup("bind-address"))
	err = viperScanLocalGitRepo.BindPFlag("bind-port", scanLocalGitRepoCmd.Flags().Lookup("bind-port"))
	err = viperScanLocalGitRepo.BindPFlag("commit-depth", scanLocalGitRepoCmd.Flags().Lookup("commit-depth"))
	err = viperScanLocalGitRepo.BindPFlag("csv", scanLocalGitRepoCmd.Flags().Lookup("csv"))
	err = viperScanLocalGitRepo.BindPFlag("debug", scanLocalGitRepoCmd.Flags().Lookup("debug"))
	err = viperScanLocalGitRepo.BindPFlag("hide-secrets", scanLocalGitRepoCmd.Flags().Lookup("hide-secrets"))
	err = viperScanLocalGitRepo.BindPFlag("ignore-extension", scanLocalGitRepoCmd.Flags().Lookup("ignore-extension"))
	err = viperScanLocalGitRepo.BindPFlag("ignore-path", scanLocalGitRepoCmd.Flags().Lookup("ignore-extension"))
	err = viperScanLocalGitRepo.BindPFlag("in-mem-clone", scanLocalGitRepoCmd.Flags().Lookup("in-mem-clone"))
	err = viperScanLocalGitRepo.BindPFlag("json", scanLocalGitRepoCmd.Flags().Lookup("json"))
	err = viperScanLocalGitRepo.BindPFlag("local-dirs", scanLocalGitRepoCmd.Flags().Lookup("local-dirs"))
	err = viperScanLocalGitRepo.BindPFlag("match-level", scanLocalGitRepoCmd.Flags().Lookup("match-level"))
	err = viperScanLocalGitRepo.BindPFlag("max-file-size", scanLocalGitRepoCmd.Flags().Lookup("max-file-size"))
	err = viperScanLocalGitRepo.BindPFlag("no-expand-orgs", scanLocalGitRepoCmd.Flags().Lookup("no-expand-orgs"))
	err = viperScanLocalGitRepo.BindPFlag("output-dir", scanGithubCmd.Flags().Lookup("output-dir"))
	err = viperScanLocalGitRepo.BindPFlag("output-prefix", scanGithubCmd.Flags().Lookup("output-prefix"))
	err = viperScanLocalGitRepo.BindPFlag("num-threads", scanLocalGitRepoCmd.Flags().Lookup("num-threads"))
	err = viperScanLocalGitRepo.BindPFlag("scan-tests", scanLocalGitRepoCmd.Flags().Lookup("scan-tests"))
	err = viperScanLocalGitRepo.BindPFlag("signature-file", scanLocalGitRepoCmd.Flags().Lookup("signature-file"))
	err = viperScanLocalGitRepo.BindPFlag("silent", scanLocalGitRepoCmd.Flags().Lookup("silent"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
