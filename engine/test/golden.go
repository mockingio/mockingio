package test

import (
	"flag"
	"io/ioutil"
	"testing"
)

var golden = flag.Bool("golden", false, "update .golden.yaml or .golden.json files")

func UpdateGoldenFile(t *testing.T, goldenFile string, data []byte) {
	t.Helper()

	if *golden {
		if err := ioutil.WriteFile(goldenFile, data, 0644); err != nil {
			t.Error(err)
		}
	}
}

func ReadGoldenFile(t *testing.T, goldenFile string) string {
	data, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Error(err)
	}

	return string(data)
}
