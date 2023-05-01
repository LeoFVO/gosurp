package cmd

import (
	"io"
	"os"

	"github.com/LeoFVO/go4hackers/cmd/hello"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "go4hackers",
	Short: "A CLI for Go4Hackers",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cliFlag, _ := cmd.Flags().GetCount("verbose")

		switch cliFlag {
			case 1:
				log.SetLevel(log.InfoLevel)
				log.Info("Log level set to INFO")
			case 2:
				log.SetLevel(log.DebugLevel)
				log.Info("Log level set to DEBUG")
			case 3:
				log.SetLevel(log.TraceLevel)
				log.Info("Log level set to TRACE")
			default:
				log.SetOutput(io.Discard)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.AddCommand(hello.RootCmd)

	RootCmd.PersistentFlags().CountP("verbose", "v", "Level of verbosity: -v for INFO, -vv for DEBUG, -vvv for TRACE.")
}