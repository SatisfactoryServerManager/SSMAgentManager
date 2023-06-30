package agents

import (
	"log"

	"github.com/SatisfactoryServerManager/SSMAgentManager/agent"
	"github.com/SatisfactoryServerManager/SSMAgentManager/gui"
	"github.com/spf13/cobra"
)

var createCmdNameFlag string
var createCmdTypeFlag string
var createCmdPortOffsetFlag int
var createCmdMemoryFlag int
var createCmdDataDirFlag string

func init() {
	Cmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new ssm agent",
	Long:  `Creates a new ssm agent and adds it to your SSM account`,
	Run: func(cmd *cobra.Command, args []string) {
		agent.LoadAgents(gui.MainApp.Preferences())
		_, err := agent.CreateNewAgent(
			createCmdNameFlag,
			createCmdTypeFlag,
			createCmdPortOffsetFlag,
			createCmdMemoryFlag,
			createCmdDataDirFlag,
			gui.MainApp.Preferences(),
		)

		if err != nil {
			log.Printf("Error creating agent, with error %s\r\n", err.Error())
			return
		}
	},
}

func init() {
	createCmd.Flags().StringVarP(&createCmdNameFlag, "name", "n", "", "The SSM Agent Name")
	createCmd.Flags().StringVarP(&createCmdTypeFlag, "type", "t", "docker", "The SSM Agent Type [docker|standalone]")
	createCmd.Flags().IntVarP(&createCmdPortOffsetFlag, "portoffset", "p", 0, "The SSM Agent Port Offset")
	createCmd.Flags().IntVarP(&createCmdMemoryFlag, "memory", "m", 0, "The SSM Agent Docker Memory Limit")
	createCmd.Flags().StringVarP(&createCmdDataDirFlag, "datadir", "d", "", "The SSM Agent Standalone Data Directory")

	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagDirname("datadir")
}
