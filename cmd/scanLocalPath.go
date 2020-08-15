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

var viperScanLocalPath *viper.Viper

// scanLocalPathCmd represents the scanLocalFiles command
var scanLocalPathCmd = &cobra.Command{
	Use:   "scanLocalPath",
	Short: "Scan local files and directorys",
	Long:  "Scan local files and directorys",
	Run: func(cmd *cobra.Command, args []string) {

		scanType := "localPath"
		sess, err := core.NewSession(viperScanLocalPath, scanType)

		_ = core.CheckArgs(sess.LocalFiles, sess.LocalDirs, sess) //YELLOW

		// exclude the .git directory from local scans as it is not handled properly here
		sess.SkippablePath = core.AppendIfMissing(sess.SkippablePath, ".git/") //YELLOW

		if err != nil {
			sess.Out.Error(err.Error()) //YELLOW
			os.Exit(1)
		}

		//sess.Out.Info("%s\n\n", common.ASCIIBanner)
		sess.Out.Important("%s v%s started at %s\n", core.Name, version.AppVersion(), sess.Stats.StartedAt.Format(time.RFC3339))
		sess.Out.Important("Loaded %d signatures.\n", len(core.Signatures))
		sess.Out.Important("Web interface available at http://%s:%d\n", sess.BindAddress, sess.BindPort)

		// Run either a file scan directly, or if it is a directory then walk the path and gather eligible files and then run a scan against each of them
		for _, fl := range sess.LocalFiles {
			if fl != "" {
				if !core.PathExists(fl, sess) {
					sess.Out.Error("\n[*] <%s> does not exist! Quitting.\n", fl)
				} else {
					core.DoFileScan(fl, sess)
				}
			}
		}

		for _, pth := range sess.LocalDirs {
			if pth != "" {
				if !core.PathExists(pth, sess) {
					sess.Out.Error("\n[*] <%s> does not exist! Quitting.\n", pth)
				} else {
					core.ScanDir(pth, sess)
				}
			}
		}

		sess.Finish()

		core.PrintSessionStats(sess)

		if !sess.Silent {
			sess.Out.Important("Press Ctrl+C to stop web server and exit.")
			select {}
		}

	},
}

func init() {
	rootCmd.AddCommand(scanLocalPathCmd)

	viperScanLocalPath = core.SetConfig()

	scanLocalPathCmd.Flags().Bool("debug", false, "Print debugging information")
	scanLocalPathCmd.Flags().Bool("hide-secrets", false, "Show secrets in any supported output")
	scanLocalPathCmd.Flags().Bool("scan-tests", false, "Scan suspected test files")
	scanLocalPathCmd.Flags().Bool("silent", false, "Suppress all output except for errors")
	scanLocalPathCmd.Flags().Int64("max-file-size", 50, "Max file size to scan")
	scanLocalPathCmd.Flags().Int("match-level", 3, "The match level of the expressions used to find matches")
	scanLocalPathCmd.Flags().String("ignore-extension", "", "a list of extensions to ignore during a scan")
	scanLocalPathCmd.Flags().String("ignore-path", "", "a list of paths to ignore during a scan")
	scanLocalPathCmd.Flags().String("rules-file", "$HOME/grover/rules/default.yml", "file(s) containing secrets detection rules.")
	scanLocalPathCmd.Flags().String("scan-dir", "", "scan a directory of files not from a git project")
	scanLocalPathCmd.Flags().String("scan-file", "", "scan a single file")

	viperScanLocalPath.BindPFlag("debug", scanLocalPathCmd.Flags().Lookup("debug"))                       //ORANGE
	viperScanLocalPath.BindPFlag("hide-secrets", scanLocalPathCmd.Flags().Lookup("hide-secrets"))         //ORANGE
	viperScanLocalPath.BindPFlag("scan-tests", scanLocalPathCmd.Flags().Lookup("scan-tests"))             //ORANGE
	viperScanLocalPath.BindPFlag("silent", scanLocalPathCmd.Flags().Lookup("silent"))                     //ORANGE
	viperScanLocalPath.BindPFlag("max-file-size", scanLocalPathCmd.Flags().Lookup("max-file-size"))       //ORANGE
	viperScanLocalPath.BindPFlag("match-level", scanLocalPathCmd.Flags().Lookup("match-level"))           //ORANGE
	viperScanLocalPath.BindPFlag("ignore-extension", scanLocalPathCmd.Flags().Lookup("ignore-extension")) //ORANGE
	viperScanLocalPath.BindPFlag("ignore-path", scanLocalPathCmd.Flags().Lookup("ignore-path"))           //ORANGE
	viperScanLocalPath.BindPFlag("rules-file", scanLocalPathCmd.Flags().Lookup("rules-file"))             //ORANGE
	viperScanLocalPath.BindPFlag("scan-dir", scanLocalPathCmd.Flags().Lookup("scan-dir"))                 //ORANGE
	viperScanLocalPath.BindPFlag("scan-file", scanLocalPathCmd.Flags().Lookup("scan-file"))               //ORANGE

}
