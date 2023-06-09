package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use: "anymon",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use anymon --help for more information")
	},
}

func init() {
	rootCmd.AddCommand(createProjectCmd)
	rootCmd.AddCommand(listenCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
