package cmd

import (
	"fmt"
	"os"

	"github.com/dpliakos/jorge/internal/jorge"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Selects or creates an environment",
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		newEnv, _ := cmd.Flags().GetBool("new")

		if debug {
			log.SetLevel(log.DebugLevel)
		}

		var selectedEnv string

		if len(args) > 0 {
			selectedEnv = args[0]
		} else {
			selectedEnv = "default"
		}

		bytes, err := jorge.UseConfigFile(selectedEnv, newEnv)

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		} else if bytes < 0 {
			fmt.Fprintln(os.Stderr, "Could not use the target file")
			os.Exit(int(bytes))
		} else if bytes == 0 {
			fmt.Fprintln(os.Stderr, "Target file is empty. Nothing to do")
		} else if bytes > 0 {
			fmt.Println(fmt.Sprintf("Using environment %s", selectedEnv))
		} else {
			fmt.Fprintln(os.Stderr, "Undefined error")
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
	useCmd.Flags().BoolP("new", "n", false, "Create a new environment")
}
