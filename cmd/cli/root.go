package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mockingio/mockingio/cmd/version"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mockingio",
	Short: "mockingio command",
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
		reportError(err)
	}
}

func reportError(err error) {
	if val, _ := os.LookupEnv("MOCKINGIO_DEBUG"); val == "1" {
		panic(err)
	}
	_, _ = fmt.Fprintf(os.Stderr, "%s", err)
	os.Exit(1)
}
