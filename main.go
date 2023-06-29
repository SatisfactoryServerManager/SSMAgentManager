package main

import (
	"github.com/SatisfactoryServerManager/SSMAgentManager/agent"
	"github.com/SatisfactoryServerManager/SSMAgentManager/cmd"
	"github.com/SatisfactoryServerManager/SSMAgentManager/gui"
)

func main() {

	gui.Init()
	agent.LoadAgents(gui.MainApp.Preferences())
	cmd.Execute()

}
