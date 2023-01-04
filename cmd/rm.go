package cmd

import (
	"fmt"
	"os"

	"github.com/dpliakos/jorge/internal/jorge"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rmCmd represents the ls command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove an environment",
	Long: `Removes a configuration environment.
	Usage:

	jorge rm <env_name>`,
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")

		if debug {
			log.SetLevel(log.DebugLevel)
		}

		var selectedEnv string
		if len(args) > 0 {
			selectedEnv = args[0]
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", "No environment specified")
			os.Exit(1)
		}

		if err := jorge.RemoveEnv(selectedEnv); err != nil {
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
			fmt.Println("Removed environment", selectedEnv)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
