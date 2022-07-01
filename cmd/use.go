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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		bytes, err := jorge.UseConfigFile(".env", selectedEnv, newEnv)

		if err != nil {
			panic(err)
		} else if bytes < 0 {
			fmt.Println("Could not use the target file")
			os.Exit(int(bytes))
		} else if bytes == 0 {
			fmt.Println("Target file is empty. Nothing to do")
		} else if bytes > 0 {
			fmt.Println(fmt.Sprintf("Using environment %s", selectedEnv))
		} else {
			fmt.Println("What the actual fuck?")
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
	useCmd.Flags().BoolP("new", "n", false, "Create a new environment")
}
