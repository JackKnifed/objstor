package main

import (
	"fmt"
	"github.com/minio/minio-go"
	"log"
	"net/http"
	"os"
)

// var ErrBadAction = errors.New("somepkg: a bad action was performed")

var endpoint = "liquidweb.services"
var bypassEncryption = false
var timeFormat = "Jan _2 2006 15:04"

type params struct {
	accessKey string
	secretKey string
	command   string
	pwd       string
	bucket    string
	cmdParams []string
}

func getConfig() params {
	var config params
	// parameters are passed as:
	// binary Command Pwd [CmdParams ...] Bucket AccessKey
	config.command = os.Args[1]
	//Pwd := os.Args[2]
	config.bucket = os.Args[len(os.Args)-2]
	config.accessKey = os.Args[len(os.Args)-1]

	config.cmdParams = os.Args[3 : len(os.Args)-2]

	// SecretKey is passed via enviroment variable
	config.secretKey = os.Getenv("PASSWORD")

	return config
}

func main() {
	config := getConfig()

	client, err := minio.New(endpoint, config.accessKey, config.secretKey, bypassEncryption)
	if err != nil {
		log.Fatalf("Access Key [%s]\nendpoint[%s]\nfailed to create new client\n%v", config.accessKey, endpoint, err)
	}

	if err = client.BucketExists(config.bucket); err != nil {
		log.Fatalf("Access Key [%s]\nendpoint[%s]\nbucket [%s] does not exist\n%v", config.accessKey, endpoint, config.bucket, err)
	}

	switch config.command {
	case "ls":
		Lsdir(client, config)
	case "mkdir":
	case "chdir":
		Chdir(client, config)
	case "rmdir":
		rmdir(client, config)
	case "delete":
		delete(client, config)
	case "get":
		get(client, config)
	case "put":
		put(client, config)
	default:
		log.Fatal("bad action")
	}

	if err != nil {
		log.Fatal(err.Error())
	}
}

func fileContentType(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Reset the read pointer if necessary.
	file.Seek(0, 0)

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(buffer), nil
}

// removes a file at a given location
// cli: `delete` `rmdir` `Pwd` `path` `bucketName` `username`
// passed to this is ["path"]
func delete(client minio.CloudStorageClient, config params) {
	err := client.RemoveObject(config.bucket, config.cmdParams[0])
	if err != nil {
		log.Fatalf("failed to delete %s\n%v", config.cmdParams[0], err)
	}
}

// does almost nothing - not required, but must return the path
// cli: `binary` `chdir` `Pwd` `path` `bucketName` `username`
func Chdir(client minio.CloudStorageClient, config params) {
	_, err := fmt.Println(config.cmdParams[0])
	if err != nil {
		panic(fmt.Sprintf("failed to print the given path %s\n%v", config.cmdParams[9], err))
	}
}

// lists the content of a directory on the remote system
// cli: `binary` `ls` `Pwd` `path` `bucketName` `username`
// passed to this is ["path"]
func Lsdir(client minio.CloudStorageClient, config params) {
	var stop chan struct{}
	var item minio.ObjectInfo
	var err error

	more := true
	res := client.ListObjects(config.bucket, config.cmdParams[0], false, stop)

	for {
		item, more = <-res
		if !more {
			return
		}
		_, err = fmt.Printf("-rwxr-xr-1 %s %s %d %s %s", item.Owner.DisplayName, item.Owner.DisplayName, item.Size, item.LastModified.Format(timeFormat), item.Key)
		if err != nil {
			stop <- struct{}{}
			log.Fatalf("failed display the file %s\n%v", item.Key, err)
		}
	}
}

// removes everything under the given path on the remote Bucket
// cli: `binary` `rmdir` `Pwd` `path` `bucketName` `username`
// passed to this is ["path"]
func rmdir(client minio.CloudStorageClient, config params) {
	var stop chan struct{}
	var item minio.ObjectInfo
	var err error

	more := true
	res := client.ListObjects(config.bucket, config.cmdParams[0], false, stop)

	for {
		item, more = <-res
		if !more {
			return
		}
		if err = client.RemoveObject(config.bucket, item.Key); err != nil {
			stop <- struct{}{}
			log.Fatalf("emptying directory [%s]\nfailed removing the file %s\n%v", config.cmdParams[0], item.Key, err)
		}
	}
}

func get(client minio.CloudStorageClient, config params) {
	err := client.FGetObject(config.bucket, config.cmdParams[1], config.cmdParams[0])
	if err != nil {
		log.Fatalf("failed to put file type - %v", err)
	}
}

func put(client minio.CloudStorageClient, config params) {
	contentType, err := fileContentType(config.cmdParams[0])
	if err != nil {
		log.Fatalf("failed to determine content type - %v", err)
	}
	_, err = client.FPutObject(config.bucket, config.cmdParams[1], config.cmdParams[0], contentType)
	if err != nil {
		log.Fatalf("failed to put file type - %v", err)
	}
}
