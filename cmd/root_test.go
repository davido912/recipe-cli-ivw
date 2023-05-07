//go:build integration

// since this is the entrypoint to the application - the tests written here will test the whole flow

package cmd

import (
	"github.com/davido912-recipe-count-test-2020/internal/testutils"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestRun(t *testing.T) {
	testDataDirPath := path.Join(testutils.GitRoot, "testdata")
	inputFilePath := path.Join(testDataDirPath, "input.json")
	outputFile, err := os.CreateTemp("", "output.json")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(outputFile.Name()) }()

	rootCmd.SetArgs([]string{
		"--file", inputFilePath,
		"-m", "Dill",
		"-p", "10335",
		"--from", "4PM",
		"--to", "10PM",
		"-o", outputFile.Name(),
	})

	Run()

	got, err := os.ReadFile(outputFile.Name())
	if err != nil {
		panic(err)
	}
	want, err := os.ReadFile(path.Join(testDataDirPath, "output.json"))
	if err != nil {
		panic(err)
	}

	require.JSONEq(t, string(want), string(got))

}
