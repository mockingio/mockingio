package cli

import (
	"context"
	"encoding/json"
	"fmt"
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

	"github.com/smockyio/smocky/api"
	"github.com/smockyio/smocky/server"
	"github.com/tuongaz/smocky-engine/engine/mock"
	"github.com/tuongaz/smocky-engine/engine/persistent/memory"
)

var filenames []string
var adminPort = 2601
var enableAdmin = false
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
		mockFileMap := map[string]string{}

		for _, filename := range filenames {
			loadedMock, err := mock.FromFile(filename)
			mock.AddIDs(loadedMock) // generate random ids
			if err != nil {
				panic(err)
			}

			if err := db.SetMock(ctx, loadedMock); err != nil {
				panic(err)
			}
			mockFileMap[loadedMock.ID] = filename

			if err := db.SetActiveSession(ctx, loadedMock.ID, uuid.NewString()); err != nil {
				panic(err)
			}

			if _, err := server.Start(ctx, loadedMock, db); err != nil {
				fmt.Printf("Failed to start server with file %v. Error: %v\n", filename, err)
				quit()
			}
		}

		if filePersist {
			db.SubscribeMockChanges(func(mock mock.Mock) {
				if err := toFile(mock, mockFileMap[mock.ID]); err != nil {
					log.WithError(err).Errorf("failed to write mock to file %s", mockFileMap[mock.ID])
				}
			})
		}

		out := output{
			URLS: server.GetServerURLs(),
		}

		if enableAdmin {
			apiServ := api.NewServer(db)
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
	startCmd.Flags().BoolVar(&enableAdmin, "admin", false, "start with admin")
	startCmd.Flags().BoolVar(&filePersist, "persist", false, "save changes to files")
	_ = startCmd.MarkFlagRequired("filename")
}
