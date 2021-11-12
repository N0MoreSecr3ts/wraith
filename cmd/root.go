// Package cmd represents the specific commands that the user will execute. Only specific code related to the command
// should be in these files. As much of the code as possible should be pushed to other packages.
package cmd

import (
	"fmt"
	"os"

	"github.com/N0MoreSecr3ts/wraith/core"
	"github.com/N0MoreSecr3ts/wraith/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (

	// wraithConfig holds the configuration data the commands
	wraithConfig *viper.Viper

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "wraith",
		Short: "A tool to scan for secrets in various digital hiding spots",
		Long:  "A tool to scan for secrets in various digital hiding spots - v" + version.AppVersion(), // TODO write a better long description
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	wraithConfig = core.SetConfig()

	rootCmd.PersistentFlags().String("bind-address", "127.0.0.1", "The IP address for the webserver")
	rootCmd.PersistentFlags().Int("bind-port", 9393, "The port for the webserver")
	rootCmd.PersistentFlags().Int("confidence-level", 3, "The confidence level level of the expressions used to find matches")
	rootCmd.PersistentFlags().Bool("csv", false, "output csv format")
	rootCmd.PersistentFlags().Bool("debug", false, "Print available debugging information to stdout")
	rootCmd.PersistentFlags().Bool("hide-secrets", false, "Do not print secrets to any supported output")
	rootCmd.PersistentFlags().StringSlice("ignore-extension", nil, "List of file extensions to ignore")
	rootCmd.PersistentFlags().StringSlice("ignore-path", nil, "List of file paths to ignore")
	rootCmd.PersistentFlags().Bool("json", false, "output json format")
	rootCmd.PersistentFlags().Int("max-file-size", 10, "Max file size to scan (in MB)")
	rootCmd.PersistentFlags().Int("num-threads", -1, "Number of execution threads")
	rootCmd.PersistentFlags().Bool("scan-tests", false, "Scan suspected test files")
	rootCmd.PersistentFlags().String("signature-file", "$HOME/.wraith/signatures/default.yaml", "file(s) containing detection signatures.")
	rootCmd.PersistentFlags().String("signature-path", "$HOME/.wraith/signatures", "path containing detection signatures.")
	rootCmd.PersistentFlags().Bool("silent", false, "Suppress all output. An alternative output will need to be configured")
	rootCmd.PersistentFlags().Bool("web-server", false, "Enable the web interface for scan output")

	err := wraithConfig.BindPFlag("bind-address", rootCmd.PersistentFlags().Lookup("bind-address"))
	err = wraithConfig.BindPFlag("bind-port", rootCmd.PersistentFlags().Lookup("bind-port"))
	err = wraithConfig.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	err = wraithConfig.BindPFlag("confidence-level", rootCmd.PersistentFlags().Lookup("confidence-level"))
	err = wraithConfig.BindPFlag("csv", rootCmd.PersistentFlags().Lookup("csv"))
	err = wraithConfig.BindPFlag("hide-secrets", rootCmd.PersistentFlags().Lookup("hide-secrets"))
	err = wraithConfig.BindPFlag("ignore-extension", rootCmd.PersistentFlags().Lookup("ignore-extension"))
	err = wraithConfig.BindPFlag("ignore-path", rootCmd.PersistentFlags().Lookup("ignore-path"))
	err = wraithConfig.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
	err = wraithConfig.BindPFlag("max-file-size", rootCmd.PersistentFlags().Lookup("max-file-size"))
	err = wraithConfig.BindPFlag("num-threads", rootCmd.PersistentFlags().Lookup("num-threads"))
	err = wraithConfig.BindPFlag("scan-tests", rootCmd.PersistentFlags().Lookup("scan-tests"))
	err = wraithConfig.BindPFlag("signature-file", rootCmd.PersistentFlags().Lookup("signature-file"))
	err = wraithConfig.BindPFlag("signature-path", rootCmd.PersistentFlags().Lookup("signature-path"))
	err = wraithConfig.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))
	err = wraithConfig.BindPFlag("web-server", rootCmd.PersistentFlags().Lookup("web-server"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}

}
