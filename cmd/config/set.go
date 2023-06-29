package config

import (
	"github.com/SatisfactoryServerManager/SSMAgentManager/gui"
	"github.com/spf13/cobra"
)

var ssmUrlFlag string
var ssmApiKeyFlag string

func init() {

	Cmd.AddCommand(setCmd)
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Updates the manager config",
	Long:  `Updates the manager config`,
	Run: func(cmd *cobra.Command, args []string) {
		gui.MainApp.Preferences().SetString("ssmurl", ssmUrlFlag)
		gui.MainApp.Preferences().SetString("ssmapikey", ssmApiKeyFlag)
	},
}

func init() {
	setCmd.Flags().StringVarP(&ssmUrlFlag, "ssmurl", "s", "https://ssmcloud.hostxtra.co.uk", "The SSM Cloud URL")
	setCmd.Flags().StringVarP(&ssmApiKeyFlag, "ssmapikey", "a", "", "The SSM Cloud API Key")

}
