package main

import (
	//"github.com/minio/minio-go"
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
		accessKey: os.Getenv("TESTUSER"),
		secretKey: os.Getenv("TESTPASS"),
		bucket:    os.Getenv("BUCKET"),
	}
	if config.accessKey != "" && config.secretKey != "" && config.bucket != ""{
		_, err := getClient(config)
		assert.Nil(t, err, "%v", err)
	}
}

func TestExpandPath(t *testing.T) {
	testData := []struct {
		pwd      string
		path     string
		expected string
	}{
		{
			pwd:      "/path/to/file",
			path:     "filename.txt",
			expected: "/path/to/file/filename.txt",
		}, {
			pwd:      "path/to/file",
			path:     "filename.txt",
			expected: "/path/to/file/filename.txt",
		}, {
			pwd:      "/path/to/file/",
			path:     "filename.txt",
			expected: "/path/to/file/filename.txt",
		}, {
			pwd:      "/path/to/file/",
			path:     "/full/path/filename.txt",
			expected: "/full/path/filename.txt",
		}, {
			pwd:      "/path/to/file",
			path:     "filename.txt/",
			expected: "/path/to/file/filename.txt",
		}, {
			pwd:      "/path///to/file",
			path:     "filename.txt",
			expected: "/path/to/file/filename.txt",
		}, {
			pwd:      "/path/./to/file",
			path:     "filename.txt",
			expected: "/path/to/file/filename.txt",
		}, {
			pwd:      "/path/to/file",
			path:     ".",
			expected: "/path/to/file",
		}, {
			pwd:      "/path/to/file",
			path:     "",
			expected: "/path/to/file",
		}, {
			pwd:      "/path/../to/file",
			path:     "filename.txt",
			expected: "/to/file/filename.txt",
		}, {
			pwd:      "/../path/to/file",
			path:     "filename.txt",
			expected: "/path/to/file/filename.txt",
		}, {
			pwd:      "/path/to/file",
			path:     "..",
			expected: "/path/to",
		}, {
			pwd:      "",
			path:     "/",
			expected: "/",
		},
	}

	var result string
	for id, test := range testData {
		result = expandPath(test.pwd, test.path)
		assert.Equal(t, test.expected, result, "test %d pwd [%q] path [%q] - did not match", id, test.pwd, test.path)
	}
}

func TestMulti(t *testing.T){
	config := params{
		accessKey: os.Getenv("TESTUSER"),
		secretKey: os.Getenv("TESTPASS"),
		bucket:    os.Getenv("BUCKET"),
	}
	if config.accessKey == "" || config.secretKey == "" || config.bucket == ""{
		t.Skip("credentials not provided, skipping tests - user [%q] bucket [%q]", config.accessKey, config.bucket)
	}
		_, err := getClient(config)
		assert.Nil(t, err, "%v", err)



}

// func TestChdir(t *testing.T) {
// 	client, _ := minio.New(endpoint, os.Getenv("accessKey"), os.Getenv("testpass"), false)
// 	r, w := io.Pipe()

// 	go chdir(client, params{cmdParams: []string{"/test/folder/within/bucket"}}, w)
// 	r.Read()

// }

// func TestFileContentType(t *testing.T) {

// }
