package config

import (
	"github.com/SatisfactoryServerManager/SSMAgentManager/gui"
	"github.com/spf13/cobra"
)

var ssmUrlFlag string
var ssmApiKeyFlag string
var ssmUserFlag string
var ssmPasswordFlag string

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
		gui.MainApp.Preferences().SetString("ssmuser", ssmUserFlag)
		gui.MainApp.Preferences().SetString("ssmpass", ssmPasswordFlag)
	},
}

func init() {
	setCmd.Flags().StringVarP(&ssmUrlFlag, "ssmurl", "s", "https://ssmcloud.hostxtra.co.uk", "The SSM Cloud URL")
	setCmd.Flags().StringVarP(&ssmApiKeyFlag, "ssmapikey", "a", "", "The SSM Cloud API Key")

	setCmd.Flags().StringVarP(&ssmUserFlag, "ssmemail", "e", "", "The SSM Cloud Account Email")
	setCmd.Flags().StringVarP(&ssmPasswordFlag, "ssmpass", "p", "", "The SSM Cloud Account Password")

	setCmd.MarkFlagsMutuallyExclusive("ssmemail", "ssmapikey")

}
