package config

import (
	"fmt"
	"path/filepath"

	"github.com/SatisfactoryServerManager/SSMAgentManager/gui"
	"github.com/spf13/cobra"
)

func init() {

	Cmd.AddCommand(printCmd)
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Prints the current manager config",
	Run: func(cmd *cobra.Command, args []string) {
		prefs := gui.MainApp.Preferences()
		fmt.Println("Config File:", filepath.Join(gui.MainApp.Storage().RootURI().Path(), "preferences.json"))
		fmt.Println("SSM Cloud URL: ", prefs.String("ssmurl"))
		fmt.Println("SSM Cloud API Key: ", prefs.String("ssmapikey"))
		fmt.Println("SSM Cloud Account Email: ", prefs.String("ssmuser"))
		fmt.Println("SSM Cloud Account Password: ", prefs.String("ssmpass"))
	},
}
