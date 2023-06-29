package agents

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SatisfactoryServerManager/SSMAgentManager/agent"
	"github.com/SatisfactoryServerManager/SSMAgentManager/gui"
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all agents",
	Run: func(cmd *cobra.Command, args []string) {
		agent.LoadAgents(gui.MainApp.Preferences())

		b, err := json.MarshalIndent(agent.AllAgents.Agents, "", "    ")

		if err != nil {
			log.Printf("Error listing agent with error %s\r\n", err.Error())
			return
		}

		fmt.Println(string(b))
	},
}
