package cmd

import (
	"fmt"
	"github.com/kaazedev/anymon/internal/project"
	"github.com/spf13/cobra"
)

var createProjectCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Run: func(cmd *cobra.Command, args []string) {
		err := project.CreateProject()
		if err != nil {
			fmt.Println("Error creating project:", err)
		}
	},
}

var listenCmd = &cobra.Command{

	Use:   "listen",
	Short: "Listen for file changes",
	Run: func(cmd *cobra.Command, args []string) {
		curr, err := project.ParseProject()
		if err != nil {
			fmt.Println("Error listening for file changes:", err)
		}

		fmt.Println("Listening for file changes...")

		project.Watch(curr)
	},
}
