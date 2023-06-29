package cmd

import (
	"fmt"
	"os"

	"github.com/SatisfactoryServerManager/SSMAgentManager/cmd/agents"
	"github.com/SatisfactoryServerManager/SSMAgentManager/cmd/config"
	"github.com/SatisfactoryServerManager/SSMAgentManager/gui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ssmagentmanager",
	Short: "SSM Agent Manager",
	Long:  "SSM Agent Manager to manage installed SSM Agents",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			gui.SetupGUI()
		}
	},
}

func init() {
	rootCmd.AddCommand(config.Cmd)
	rootCmd.AddCommand(agents.Cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
