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

var viperScanLocalGitRepo *viper.Viper

// scanLocalGitRepoCmd represents the scanLocalGitRepo command
var scanLocalGitRepoCmd = &cobra.Command{
	Use:   "scanLocalGitRepo",
	Short: "Scan a git repo on a local machine",
	Long:  "Scan a git repo on a local machine",
	Run: func(cmd *cobra.Command, args []string) {

		scanType := "localGit"
		sess, err := core.NewSession(viperScanLocalGitRepo, scanType)

		if err != nil {
			sess.Out.Error("Failed to generate a new session: %s\n", err.Error()) //YELLOW
			os.Exit(1)
		}

		//sess.Out.Info("%s\n\n", common.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", "127.0.0.1", 9393)

		core.GatherLocalRepositories(sess)
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
	rootCmd.AddCommand(scanLocalGitRepoCmd)

	viperScanLocalGitRepo = core.SetConfig()

	scanLocalGitRepoCmd.Flags().Bool("debug", false, "Print debugging information")
	scanLocalGitRepoCmd.Flags().Bool("hide-secrets", false, "Hide secrets from output")
	scanLocalGitRepoCmd.Flags().Bool("in-mem-clone", false, "Clone repos in memory")
	scanLocalGitRepoCmd.Flags().Bool("no-expand-orgs", false, "Don't add members to targets when processing organizations")
	scanLocalGitRepoCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanLocalGitRepoCmd.Flags().Bool("silent", false, "No output")
	scanLocalGitRepoCmd.Flags().Int("bind-port", 9393, "The port for the webserver")
	scanLocalGitRepoCmd.Flags().Int("commit-depth", 0, "Set the depth for commits")
	scanLocalGitRepoCmd.Flags().Int("match-level", 3, "Signature match level")
	scanLocalGitRepoCmd.Flags().Int("max-file-size", 50, "Max file size to scan")
	scanLocalGitRepoCmd.Flags().Int("num-threads", 0, "The number of threads to execute with")
	scanLocalGitRepoCmd.Flags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	scanLocalGitRepoCmd.Flags().String("ignore-extension", "", "a comma separated list of extensions to ignore")
	scanLocalGitRepoCmd.Flags().String("ignore-path", "", "a comma separated list of paths to ignore")
	scanLocalGitRepoCmd.Flags().String("local-dirs", "", "local disk parent dir containing git repos")
	scanLocalGitRepoCmd.Flags().String("signature-file", "$HOME/.wraith/signatures/default_signatures.yml", "file(s) containing detection signatures.")

	viperScanLocalGitRepo.BindPFlag("bind-address", scanLocalGitRepoCmd.Flags().Lookup("bind-address"))         //ORANGE
	viperScanLocalGitRepo.BindPFlag("bind-port", scanLocalGitRepoCmd.Flags().Lookup("bind-port"))               //ORANGE
	viperScanLocalGitRepo.BindPFlag("commit-depth", scanLocalGitRepoCmd.Flags().Lookup("commit-depth"))         //ORANGE
	viperScanLocalGitRepo.BindPFlag("debug", scanLocalGitRepoCmd.Flags().Lookup("debug"))                       //ORANGE
	viperScanLocalGitRepo.BindPFlag("hide-secrets", scanLocalGitRepoCmd.Flags().Lookup("hide-secrets"))         //ORANGE
	viperScanLocalGitRepo.BindPFlag("ignore-extension", scanLocalGitRepoCmd.Flags().Lookup("ignore-extension")) //ORANGE
	viperScanLocalGitRepo.BindPFlag("ignore-path", scanLocalGitRepoCmd.Flags().Lookup("ignore-extension"))      //ORANGE
	viperScanLocalGitRepo.BindPFlag("in-mem-clone", scanLocalGitRepoCmd.Flags().Lookup("in-mem-clone"))         //ORANGE
	viperScanLocalGitRepo.BindPFlag("local-dirs", scanLocalGitRepoCmd.Flags().Lookup("local-dirs"))             //ORANGE
	viperScanLocalGitRepo.BindPFlag("match-level", scanLocalGitRepoCmd.Flags().Lookup("match-level"))           //ORANGE
	viperScanLocalGitRepo.BindPFlag("max-file-size", scanLocalGitRepoCmd.Flags().Lookup("max-file-size"))       //ORANGE
	viperScanLocalGitRepo.BindPFlag("no-expand-orgs", scanLocalGitRepoCmd.Flags().Lookup("no-expand-orgs"))     //ORANGE
	viperScanLocalGitRepo.BindPFlag("num-threads", scanLocalGitRepoCmd.Flags().Lookup("num-threads"))           //ORANGE
	viperScanLocalGitRepo.BindPFlag("scan-tests", scanLocalGitRepoCmd.Flags().Lookup("scan-tests"))             //ORANGE
	viperScanLocalGitRepo.BindPFlag("signature-file", scanLocalGitRepoCmd.Flags().Lookup("signature-file"))     //ORANGE
	viperScanLocalGitRepo.BindPFlag("silent", scanLocalGitRepoCmd.Flags().Lookup("silent"))                     //ORANGE

}
