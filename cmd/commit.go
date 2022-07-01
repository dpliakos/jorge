/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("commit called")

		debug, _ := cmd.Flags().GetBool("debug")

		if debug {
			log.SetLevel(log.DebugLevel)
		}

		bytes, err := jorge.StoreConfigFile(".env", "dev")

		if err != nil {
			panic(err)
		} else if bytes < 0 {
			panic("Could not store the file")
		} else if bytes == 0 {
			fmt.Println("Target file is empty. Nothing to do")
		} else if bytes > 0 {
			fmt.Println("Target file updated")
		} else {
			fmt.Println("What the actual fuck?")
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
