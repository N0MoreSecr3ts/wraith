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

		sess.Out.Warn("%s\n\n", core.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", "127.0.0.1", 9393)

		core.GatherLocalRepositories(sess)
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
	rootCmd.AddCommand(scanLocalGitRepoCmd)

	viperScanLocalGitRepo = core.SetConfig()

	scanLocalGitRepoCmd.Flags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	scanLocalGitRepoCmd.Flags().Int("bind-port", 9393, "The port for the webserver")
	scanLocalGitRepoCmd.Flags().Int("confidence-level", 3, "The confidence level level of the expressions used to find matches")
	scanLocalGitRepoCmd.Flags().Float64("commit-depth", -1, "Set the commit depth to scan")
	scanLocalGitRepoCmd.Flags().Bool("debug", false, "Print available debugging information to stdout")
	scanLocalGitRepoCmd.Flags().Bool("hide-secrets", false, "Do not print secrets to any supported output")
	scanLocalGitRepoCmd.Flags().StringSlice("ignore-extension", nil, "List of file extensions to ignore")
	scanLocalGitRepoCmd.Flags().StringSlice("ignore-path", nil, "List of file paths to ignore")
	scanLocalGitRepoCmd.Flags().Int("max-file-size", 10, "Max file size to scan (in MB)")
	scanLocalGitRepoCmd.Flags().Int("num-threads", -1, "Number of execution threads")
	scanLocalGitRepoCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanLocalGitRepoCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default.yaml", "file(s) containing detection signatures.")
	scanLocalGitRepoCmd.Flags().String("signature-path", "$HOME/.wraith/signatures", "path containing detection signatures.")
	scanLocalGitRepoCmd.Flags().Bool("silent", false, "Suppress all output except for errors")
	scanLocalGitRepoCmd.Flags().StringSlice("local-repos", nil, "List of local git repos to scan")

	err := viperScanLocalGitRepo.BindPFlag("bind-address", scanLocalGitRepoCmd.Flags().Lookup("bind-address"))
	err = viperScanLocalGitRepo.BindPFlag("bind-port", scanLocalGitRepoCmd.Flags().Lookup("bind-port"))
	err = viperScanLocalGitRepo.BindPFlag("commit-depth", scanLocalGitRepoCmd.Flags().Lookup("commit-depth"))
	err = viperScanLocalGitRepo.BindPFlag("debug", scanLocalGitRepoCmd.Flags().Lookup("debug"))
	err = viperScanLocalGitRepo.BindPFlag("hide-secrets", scanLocalGitRepoCmd.Flags().Lookup("hide-secrets"))
	err = viperScanLocalGitRepo.BindPFlag("ignore-extension", scanLocalGitRepoCmd.Flags().Lookup("ignore-extension"))
	err = viperScanLocalGitRepo.BindPFlag("ignore-path", scanLocalGitRepoCmd.Flags().Lookup("ignore-path"))
	err = viperScanLocalGitRepo.BindPFlag("confidence-level", scanLocalGitRepoCmd.Flags().Lookup("confidence-level"))
	err = viperScanLocalGitRepo.BindPFlag("max-file-size", scanLocalGitRepoCmd.Flags().Lookup("max-file-size"))
	err = viperScanLocalGitRepo.BindPFlag("num-threads", scanLocalGitRepoCmd.Flags().Lookup("num-threads"))
	err = viperScanLocalGitRepo.BindPFlag("scan-tests", scanLocalGitRepoCmd.Flags().Lookup("scan-tests"))
	err = viperScanLocalGitRepo.BindPFlag("signature-file", scanLocalGitRepoCmd.Flags().Lookup("signature-file"))
	err = viperScanLocalGitRepo.BindPFlag("signature-path", scanLocalGitRepoCmd.Flags().Lookup("signature-path"))
	err = viperScanLocalGitRepo.BindPFlag("silent", scanLocalGitRepoCmd.Flags().Lookup("silent"))
	err = viperScanLocalGitRepo.BindPFlag("local-repos", scanLocalGitRepoCmd.Flags().Lookup("local-repos"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
