package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/smockyio/smocky/engine/mock"
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

		ctx := context.Background()

		db := memory.New()
		persistent.New(db)

		for _, filename := range filenames {
			loadedMock, err := mock.FromFile(filename)
			if err != nil {
				panic(err)
			}

			if err := db.SetMock(ctx, loadedMock); err != nil {
				panic(err)
			}

			if err := db.SetActiveSession(ctx, loadedMock.ID, uuid.NewString()); err != nil {
				panic(err)
			}

			if _, err := server.Start(ctx, loadedMock); err != nil {
				fmt.Printf("Failed to start server with file %v. Error: %v\n", filename, err)
				quit()
			}
		}

		out := output{
			URLS: server.GetServerURLs(),
		}

		if enableAdmin {
			apiServ := api.NewServer()
			adminURL, shutdownServer, err := apiServ.Start(ctx, strconv.Itoa(adminPort))
			if err != nil {
				fmt.Printf("Failed to start admin server. Error: %v\n", err)
				shutdownServer()
				quit()
			}
			out.Admin = adminURL
		}

		data, _ := json.Marshal(out)
		fmt.Println(string(data))

		<-stopSignalChanel
		server.RemoveAllServers()
		fmt.Println("servers stopped")
	},
}

func quit() {
	server.RemoveAllServers()
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringArrayVarP(&filenames, "filename", "f", []string{}, "location of the mock file")
	startCmd.Flags().IntVar(&adminPort, "admin-port", 2601, "port for admin API server")
	startCmd.Flags().BoolVar(&enableAdmin, "admin", false, "start with admin")
	_ = startCmd.MarkFlagRequired("filename")
}
