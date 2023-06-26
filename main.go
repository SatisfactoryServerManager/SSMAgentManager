package main

import (
	"github.com/SatisfactoryServerManager/SSMAgentManager/cmd"
	"github.com/SatisfactoryServerManager/SSMAgentManager/gui"
)

func main() {

	gui.Init()
	cmd.Execute()

}
