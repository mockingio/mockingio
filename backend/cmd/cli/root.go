package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	buildVersion  string
	buildRevision string
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

		fmt.Printf(
			"Version: %v, Revision: %v, Build Date: %v\n",
			buildVersion,
			buildRevision,
			time.Now().UTC().Format(time.RFC3339),
		)
	},
}

var version bool

func init() {
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "show version")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
