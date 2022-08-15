package cmd

import (
	"fmt"
	"os"

	"github.com/dpliakos/jorge/internal/jorge"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a jorge environment",
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		configFilePath, _ := cmd.Flags().GetString("config")

		if debug {
			log.SetLevel(log.DebugLevel)
		}

		err := jorge.Init(configFilePath)

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
			fmt.Println("Created new jorge project")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("config", "c", "", "Declare the project's config file path")
}
