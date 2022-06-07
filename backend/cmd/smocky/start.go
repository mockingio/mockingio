package smocky

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/smockyio/smocky/backend/server"
)

var filenames []string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a mock server",
	Long: `
smocky start --filename mock.yml
smocky start --filename mock1.yml --filename mock2.yml
`,
	Run: func(cmd *cobra.Command, args []string) {
		stopSignalChanel := make(chan os.Signal, 1)
		signal.Notify(
			stopSignalChanel,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)

		var shutdownServers []func()

		for _, filename := range filenames {
			serv := server.New()
			url, shutdownServer, err := serv.StartFromFile(context.Background(), filename)

			if err != nil {
				fmt.Printf("Failed to start server with file %v. Error: %v\n", filename, err)
				for _, shutdown := range shutdownServers {
					shutdown()
				}
				os.Exit(1)
			}
			shutdownServers = append(shutdownServers, shutdownServer)

			fmt.Printf("serving server from: %v\n", url)
		}

		<-stopSignalChanel
		for _, shutdown := range shutdownServers {
			shutdown()
		}

		fmt.Println("servers stopped")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringArrayVarP(&filenames, "filename", "f", []string{}, "mock files")
	_ = startCmd.MarkFlagRequired("filename")
}
