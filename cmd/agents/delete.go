package agents

import (
	"log"

	"github.com/SatisfactoryServerManager/SSMAgentManager/agent"
	"github.com/SatisfactoryServerManager/SSMAgentManager/gui"
	"github.com/spf13/cobra"
)

var deleteCmdNameFlag string

func init() {
	Cmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a ssm agent",
	Long:  `Deletes a new ssm agent and removes it from your SSM account`,
	Run: func(cmd *cobra.Command, args []string) {
		agent.LoadAgents(gui.MainApp.Preferences())
		err := agent.DeleteAgent(
			deleteCmdNameFlag,
			gui.MainApp.Preferences(),
		)

		if err != nil {
			log.Printf("Error deleting agent, with error %s\r\n", err.Error())
			return
		}
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&deleteCmdNameFlag, "name", "n", "", "The SSM Agent Name")

	deleteCmd.MarkFlagRequired("name")
}
