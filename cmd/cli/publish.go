package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/mockingio/mockingio/engine/mock"
)

var filename string
var endpoint string

// publishCmd represents the push command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish your mocks to mocking.io",
	Long: `
mockingio publish --filename mock.yml
`,
	Run: func(cmd *cobra.Command, args []string) {
		absFilePath, err := filepath.Abs(filename)
		if err != nil {
			reportError(err)
		}

		fileContent, err := os.ReadFile(absFilePath)
		if err != nil {
			reportError(err)
		}

		loadedMock, err := mock.FromYaml(string(fileContent))
		if err != nil {
			reportError(err)
		}

		if err := loadedMock.ApplyDefault().Validate(); err != nil {
			reportError(err)
		}

		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		req, err := http.NewRequest("POST", endpoint+"/mocks", strings.NewReader(string(fileContent)))
		if err != nil {
			reportError(err)
		}

		req.Header.Set("Content-Type", "application/yaml")
		resp, err := client.Do(req)
		if err != nil {
			reportError(err)
		}

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			reportError(fmt.Errorf("failed to publish mock. Response code: %v", resp.StatusCode))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			reportError(err)
		}

		var obj map[string]interface{}
		if err = json.Unmarshal(body, &obj); err != nil {
			reportError(err)
		}

		fmt.Println("Mock URLs")
		fmt.Printf("http://%v\n", obj["url"])
		fmt.Printf("https://%v\n", obj["url"])
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVarP(&filename, "filename", "f", "", "location of the mock file")
	publishCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "https://api.mocking.io", "location of the mock file")
	_ = publishCmd.MarkFlagRequired("filename")
}
