package agents

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "agents",
	Short: "SSM Agent Manager Agents",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Agents")
	},
}
