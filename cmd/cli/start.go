package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mockingio/engine/persistent"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/mockingio/engine/mock"
	"github.com/mockingio/engine/persistent/memory"
	"github.com/mockingio/mockingio/api"
	"github.com/mockingio/mockingio/server"
)

var filenames []string
var adminPort = 2601
var filePersist = false

type output struct {
	URLS  []string `json:"urls,omitempty"`
	Admin string   `json:"admin,omitempty"`
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a mock server",
	Long: `
mockingio start --filename mock.yml
mockingio start --filename mock1.yml --filename mock2.yml
mockingio start --filename mock.yml --output-json
`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		db := memory.New()
		mockFileMap := mustLoadMocks(ctx, filenames, db)

		serv := server.New(db)

		for _, item := range mockFileMap {
			if _, err := serv.Start(ctx, item.mock); err != nil {
				fmt.Printf("Failed to start server with file %v. Error: %v\n", item.filename, err)
				quit()
			}
		}

		if filePersist {
			db.SubscribeMockChanges(func(mock mock.Mock) {
				if err := toFile(mock, mockFileMap[mock.ID].filename); err != nil {
					log.WithError(err).Errorf("failed to write mock to file %s", mockFileMap[mock.ID].filename)
				}
			})
		}

		out := output{
			URLS: server.GetServerURLs(),
		}

		apiServ := api.NewServer(db)
		adminURL, shutdownServer, err := apiServ.Start(ctx, strconv.Itoa(adminPort))
		if err != nil {
			fmt.Printf("Failed to start admin server. Error: %v\n", err)
			shutdownServer()
			quit()
		}
		out.Admin = adminURL

		data, _ := json.Marshal(out)
		fmt.Println(string(data))

		stopSignalChanel := registerStopSignal()
		<-stopSignalChanel
		server.RemoveAllServers()
		fmt.Println("servers stopped")
	},
}

// mustLoadMocks loop through mock files and load them to database.
func mustLoadMocks(ctx context.Context, filenames []string, db persistent.Persistent) map[string]struct {
	filename string
	mock     *mock.Mock
} {
	mockFileMap := map[string]struct {
		filename string
		mock     *mock.Mock
	}{}

	for _, filename := range filenames {
		loadedMock, err := mock.FromFile(filename, mock.WithIDGeneration())
		if err != nil {
			panic(err)
		}

		if err := db.SetMock(ctx, loadedMock); err != nil {
			panic(err)
		}

		mockFileMap[loadedMock.ID] = struct {
			filename string
			mock     *mock.Mock
		}{
			filename: filename,
			mock:     loadedMock,
		}

		if err := db.SetActiveSession(ctx, loadedMock.ID, uuid.NewString()); err != nil {
			panic(err)
		}
	}

	return mockFileMap
}

func registerStopSignal() chan os.Signal {
	stopSignalChanel := make(chan os.Signal, 1)
	signal.Notify(
		stopSignalChanel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	return stopSignalChanel
}

func toFile(mock mock.Mock, filename string) error {
	text, err := yaml.Marshal(mock)
	if err != nil {
		return errors.Wrap(err, "marshal mock to yaml")
	}

	fileStats, err := os.Stat(filename)
	if err != nil {
		return errors.Wrap(err, "get file stats")
	}

	if err := ioutil.WriteFile(filename, text, fileStats.Mode().Perm()); err != nil {
		return errors.Wrap(err, "write file")
	}

	return nil
}

func quit() {
	server.RemoveAllServers()
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringArrayVarP(&filenames, "filename", "f", []string{}, "location of the mock file")
	startCmd.Flags().IntVar(&adminPort, "admin-port", 2601, "port for admin API server")
	startCmd.Flags().BoolVar(&filePersist, "persist", false, "save changes to files")
	_ = startCmd.MarkFlagRequired("filename")
}
