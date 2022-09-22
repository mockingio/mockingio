package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/mockingio/mockingio/api"
	"github.com/mockingio/mockingio/engine/database"
	"github.com/mockingio/mockingio/engine/database/memory"
	"github.com/mockingio/mockingio/engine/mock"
	"github.com/mockingio/mockingio/engine/server"
)

var filenames []string
var adminPort = 2601
var filePersist = false

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "NewMockServer a mock server",
	Long: `
mockingio start --filename mock.yml
mockingio start --filename mock1.yml --filename mock2.yml
mockingio start --filename mock.yml --output-json
`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		db := memory.New()
		mockServer := server.New(db)
		mockFileMap := mustLoadMocks(ctx, filenames, db)

		// start mock servers
		for _, item := range mockFileMap {
			if _, err := mockServer.NewMockServer(ctx, item.mock); err != nil {
				fmt.Printf("Failed to start server with file %v. Error: %v\n", item.filename, err)
				quit(mockServer)
			}
		}

		// start admin server
		adminURL, shutdownServer, err := api.NewServer(db, mockServer).Start(ctx, strconv.Itoa(adminPort))
		if err != nil {
			log.WithError(err).Error("Failed to start api server")
			shutdownServer()
			quit(mockServer)
		}

		// save mock changes to files
		if filePersist {
			db.SubscribeMockChanges(func(mock mock.Mock) {
				if err := toFile(mock, mockFileMap[mock.ID].filename); err != nil {
					log.WithError(err).Errorf("failed to write mock to file %s", mockFileMap[mock.ID].filename)
				}
			})
		}

		printServersInfo(mockServer.GetMockServerURLs(), adminURL)
		onStopSignal(mockServer.StopAllServers)
	},
}

func printServersInfo(mockUrls []string, adminURL string) {
	data, _ := json.Marshal(map[string]any{
		"urls":      mockUrls,
		"admin_url": adminURL,
	})
	fmt.Println(string(data))
}

// mustLoadMocks loop through mock files and load them to database.
func mustLoadMocks(ctx context.Context, filenames []string, db database.EngineDB) map[string]struct {
	filename string
	mock     *mock.Mock
} {
	mockFileMap := map[string]struct {
		filename string
		mock     *mock.Mock
	}{}

	for _, filename := range filenames {
		absFilePath, err := filepath.Abs(filename)
		if err != nil {
			panic(err)
		}

		loadedMock, err := mock.FromFile(absFilePath, mock.WithIDGeneration())
		if err != nil {
			if val, _ := os.LookupEnv("MOCKINGIO_DEBUG"); val == "1" {
				panic(err)
			}
			_, _ = fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
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

// onStopSignal registers a signal handler for SIGINT and SIGTERM.
func onStopSignal(handler func()) {
	stopSignalChanel := make(chan os.Signal, 1)
	signal.Notify(
		stopSignalChanel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	<-stopSignalChanel
	handler()
	log.Info("servers stopped")
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

func quit(mockServer *server.Server) {
	mockServer.StopAllServers()
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringArrayVarP(&filenames, "filename", "f", []string{}, "location of the mock file")
	startCmd.Flags().IntVar(&adminPort, "admin-port", 2601, "port for admin API server")
	startCmd.Flags().BoolVar(&filePersist, "persist", false, "save changes to files")
	_ = startCmd.MarkFlagRequired("filename")
}
