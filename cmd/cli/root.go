package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/smockyio/smocky/cmd/version"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "smocky",
	Short: "Smocky command",
	Run: func(cmd *cobra.Command, args []string) {
		showVersion, _ := cmd.Flags().GetBool("version")
		if !showVersion {
			_ = cmd.Help()
			return
		}

		fmt.Println(version.Long())
	},
}

var flagVersion bool

func init() {
	rootCmd.Flags().BoolVarP(&flagVersion, "version", "v", false, "show version")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
