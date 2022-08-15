package cmd

import (
	"fmt"
	"os"

	"github.com/dpliakos/jorge/internal/jorge"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Stores the current config file",
	Long:  `Updates the environment configuration with the current version of the configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")

		if debug {
			log.SetLevel(log.DebugLevel)
		}

		err := jorge.CommitCurrentEnv()

		if err != nil {
			if debug && err.OriginalErr != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.OriginalErr.Error())
			}

			fmt.Fprintf(os.Stderr, "%s\n", err.Message)
			fmt.Fprintf(os.Stderr, "%s\n", err.Solution)

			if err.Code > 0 {
				os.Exit(err.Code)
			} else {
				os.Exit(1)
			}
		} else {
			fmt.Println("Env committed")
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
