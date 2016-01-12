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

func TestGetClient(t *testing.T) {
	config := params{
		accessKey: os.Getenv("testuser"),
		secretKey: os.Getenv("testpass"),
		bucket:    os.Getenv("bucket"),
	}
	assert.NotEmpty(t, config.accessKey, "testuser not set for test")
	assert.NotEmpty(t, config.secretKey, "testpass not set for test")
	assert.NotEmtpy(t, config.bucket, "bucket not set for test")

	_, err := getClient(config)
	assert.Nil(t, err, "Failed to create bucket - %v", err)
}

// func TestChdir(t *testing.T) {
// 	client, _ := minio.New(endpoint, os.Getenv("accessKey"), os.Getenv("testpass"), false)
// 	r, w := io.Pipe()

// 	go chdir(client, params{cmdParams: []string{"/test/folder/within/bucket"}}, w)
// 	r.Read()

// }

// func TestFileContentType(t *testing.T) {

// }
