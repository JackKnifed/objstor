package main

import (
	"github.com/minio/minio-go"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	// parameters are passed as:
	// binary Command Pwd [CmdParams ...] Bucket AccessKey
	testData := []struct {
		input          []string
		inputSecretKey string
		output         params
	}{
		{
			input: []string{
				"null",
				"comMand",
				"PWD",
				"cmdParams",
				"bukkit",
				"aCceSsKey"},
			inputSecretKey: "seKreTkEy",
			output: params{
				accessKey: "aCceSsKey",
				secretKey: "seKreTkEy",
				command:   "comMand",
				pwd:       "PWD",
				bucket:    "bukkit",
				cmdParams: []string{
					"cmdParams",
				},
			},
		}, {
			input: []string{
				"",
				"comMand",
				"PWD",
				"cmdParams",
				"secondParam",
				"bukkit",
				"aCceSsKey"},
			inputSecretKey: "seKreTkEy",
			output: params{
				accessKey: "aCceSsKey",
				secretKey: "seKreTkEy",
				command:   "comMand",
				pwd:       "PWD",
				bucket:    "bukkit",
				cmdParams: []string{
					"cmdParams",
					"secondParam",
				},
			},
		},
	}

	for _, test := range testData {
		os.Setenv("PASSWORD", test.inputSecretKey)
		config := getConfig(test.input)
		assert.Equal(t, test.output, config)
	}
}

func ExampleChDir() {
	client, _ := minio.New("liquidweb.services", "accessKey", "sekretKey", false)
	chdir(client,
		params{
			cmdParams: []string{
				"/test/folder/within/bucket",
			},
		})

	//Output:
	// /test/folder/within/bucket
}

// func TestFileContentType(t *testing.T) {

// }
