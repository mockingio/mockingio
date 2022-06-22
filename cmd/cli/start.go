package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/smockyio/smocky/engine/mock"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/smockyio/smocky/engine/persistent"
	"github.com/smockyio/smocky/engine/persistent/memory"
	"github.com/smockyio/smocky/server"
	"github.com/smockyio/smocky/server/api"
)

var filenames []string
var adminPort = 2601
var enableAdmin = false

type output struct {
	URLS  []string `json:"urls,omitempty"`
	Admin string   `json:"admin,omitempty"`
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a mock server",
	Long: `
smocky start --filename mock.yml
smocky start --filename mock1.yml --filename mock2.yml
smocky start --filename mock.yml --output-json
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

		out := output{}
		ctx := context.Background()

		// initialise database
		persistent.New(memory.New())
		db := persistent.GetDefault()

		for _, filename := range filenames {
			// TODO: check the file extension, support loading mock from JSON
			loadedMock, err := mock.FromYamlFile(filename)
			if err != nil {
				panic(err)
			}

			if err := db.SetMock(ctx, loadedMock); err != nil {
				panic(err)
			}

			if err := db.SetActiveSession(ctx, loadedMock.ID, uuid.NewString()); err != nil {
				panic(err)
			}
			serv := backend.New()
			url, shutdownServer, err := serv.Start(ctx, loadedMock.ID)

			if err != nil {
				fmt.Printf("Failed to start server with file %v. Error: %v\n", filename, err)
				quit(shutdownServers)
			}
			shutdownServers = append(shutdownServers, shutdownServer)

			out.URLS = append(out.URLS, url)
		}

		if enableAdmin {
			apiServ := api.NewServer()
			adminURL, shutdownServer, err := apiServ.Start(ctx, strconv.Itoa(adminPort))
			if err != nil {
				fmt.Printf("Failed to start admin server. Error: %v\n", err)
				quit(shutdownServers)
			}
			shutdownServers = append(shutdownServers, shutdownServer)
			out.Admin = adminURL
		}

		data, _ := json.Marshal(out)
		fmt.Println(string(data))

		<-stopSignalChanel
		for _, shutdown := range shutdownServers {
			shutdown()
		}

		fmt.Println("servers stopped")
	},
}

func quit(shutdowns []func()) {
	for _, shutdown := range shutdowns {
		shutdown()
	}
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringArrayVarP(&filenames, "filename", "f", []string{}, "location of the mock file")
	startCmd.Flags().IntVar(&adminPort, "admin-port", 2601, "port for admin API server")
	startCmd.Flags().BoolVar(&enableAdmin, "admin", false, "start with admin")
	_ = startCmd.MarkFlagRequired("filename")
}
