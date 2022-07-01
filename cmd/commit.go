package cmd

import (
	"fmt"

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
			panic(err)
		} else {
			fmt.Println("Env committed")
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
