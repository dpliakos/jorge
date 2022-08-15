package cmd

import (
	"fmt"
	"os"

	"github.com/dpliakos/jorge/internal/jorge"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List the available environments",
	Long: `Shows a list with the configuration environments that the user created.
	Usage:

	jorge ls`,
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")

		if debug {
			log.SetLevel(log.DebugLevel)
		}

		if err := jorge.ListEnvironments(); err != nil {
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
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
