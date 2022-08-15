package cmd

import (
	"fmt"
	"os"

	"github.com/dpliakos/jorge/internal/jorge"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restores the current configuration file with the copy that is saved in the .jorge dir",
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")

		if debug {
			log.SetLevel(log.DebugLevel)
		}

		err := jorge.RestoreEnv()

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
			fmt.Println("Env restored")
		}
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
