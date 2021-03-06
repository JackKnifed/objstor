package main

import (
	"bytes"
	"fmt"
	"github.com/minio/minio-go"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// example call
//  /root/bin/objstor put / /var/cpanel/backups/tmp.test.2266-1452758896 lol-secure/validate.tmp-2266-1452758896.txt objects.liquidweb.services URHR6M3LZM2835X9BNYF

const (
	endpoint         = "objects.liquidweb.services"
	bypassEncryption = false
	folderContentType = "application/x-directory"
	timeFormat       = "Jan _2 2006 15:04"
	relPathSeperator = "/"
)

type params struct {
	accessKey string
	secretKey string
	command   string
	pwd       string
	bucket    string
	cmdParams []string
}

func expandPath(pwd, path string) string {
	var pathParts []string
	if strings.HasPrefix(path, relPathSeperator) {
		pathParts = strings.Split(path, relPathSeperator)
	} else {
		pathParts = strings.Split(pwd+relPathSeperator+path, relPathSeperator)
	}

	for i := 0; i < len(pathParts); {
		switch pathParts[i] {
		case ".":
			if i+1 < len(pathParts) {
				// if there's at least one element after this one, cut this element out
				pathParts = append(pathParts[:i], pathParts[i+1:]...)
			} else {
				// otherwise cut this element off the end
				pathParts = pathParts[:i]
			}
		case "":
			if i+1 < len(pathParts) {
				pathParts = append(pathParts[:i], pathParts[i+1:]...)
			} else {
				pathParts = pathParts[:i]
			}
		case "..":
			if i < 1 {
				// if this is the first element in the list, start at the second element
				pathParts = pathParts[i+1:]
			} else if i+1 < len(pathParts) {
				// if there is at least one more item in the list cut this element and the one before, and tack on what's left
				pathParts = append(pathParts[:i-1], pathParts[i+1:]...)
			} else {
				// else just lop off this element and the one before
				pathParts = pathParts[:i-1]
			}
		default:
			i++
		}
	}

	return relPathSeperator + strings.Join(pathParts, relPathSeperator)
}

// returns the contentType of the given file at the given location
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

// parses the config out of the Args (and one env variable)
func getConfig(args []string) params {
	var config params
	// parameters are passed as:
	// binary Command Pwd [CmdParams ...] Bucket AccessKey
	config.command = args[1]
	config.pwd = args[2]

	// working from the back
	config.bucket = args[len(args)-2]
	config.accessKey = args[len(args)-1]

	// everything in the middle is the cmdParams
	config.cmdParams = args[3 : len(args)-2]

	// SecretKey is passed via enviroment variable
	config.secretKey = os.Getenv("PASSWORD")

	return config
}

func getClient(config params) (minio.CloudStorageClient, error) {
	client, err := minio.NewV2(endpoint, config.accessKey, config.secretKey, bypassEncryption)
	if err != nil {
		return nil, fmt.Errorf("Access Key [%s]\nendpoint[%s]\nfailed to create new client\n%v", config.accessKey, endpoint, err)
	}

	if err = client.BucketExists(config.bucket); err != nil {
		//return nil, fmt.Errorf("Access Key [%s]\nendpoint[%s]\nbucket [%s] does not exist\n%v", config.accessKey, endpoint, config.bucket, err)
		return nil, err
	}
	return client, nil
}

// removes a file at a given location
// cli: `delete` `rmdir` `Pwd` `path` `bucketName` `username`
// passed to this is ["path"]
func delete(client minio.CloudStorageClient, config params) error {
	err := client.RemoveObject(config.bucket, expandPath(config.pwd, config.cmdParams[0]))
	if err != nil {
		return fmt.Errorf("failed to delete %s\n%v", config.cmdParams[0], err)
	}
	return nil
}

// does almost nothing - not required, but must return the path
// cli: `binary` `chdir` `Pwd` `path` `bucketName` `username`
func chdir(client minio.CloudStorageClient, config params, out io.Writer) error {
	_, err := fmt.Fprintln(out, expandPath(config.pwd, config.cmdParams[0]))
	if err != nil {
		return fmt.Errorf("failed to print the given path %s\n%v", config.cmdParams[9], err)
	}
	return nil
}

// lists the content of a directory on the remote system
// cli: `binary` `ls` `Pwd` `path` `bucketName` `username`
// passed to this is ["path"]
func lsdir(client minio.CloudStorageClient, config params, out io.Writer) error {
	var stop chan struct{}
	var item minio.ObjectInfo
	var err error

	more := true
	res := client.ListObjects(config.bucket, expandPath(config.pwd, config.cmdParams[0]), false, stop)

	for {
		item, more = <-res
		if !more {
			return nil
		}
		_, err = fmt.Fprintf(out, "-rwxr-xr-x %s %s %d %s %s", item.Owner.DisplayName, item.Owner.DisplayName, item.Size, item.LastModified.Format(timeFormat), item.Key)
		if err != nil {
			stop <- struct{}{}
			return fmt.Errorf("failed display the file %s\n%v", item.Key, err)
		}
	}
	return nil
}

// removes everything under the given path on the remote Bucket
// cli: `binary` `rmdir` `Pwd` `path` `bucketName` `username`
// passed to this is ["path"]
func rmdir(client minio.CloudStorageClient, config params) error {
	var stop chan struct{}
	var item minio.ObjectInfo
	var err error

	more := true
	res := client.ListObjects(config.bucket, expandPath(config.pwd, config.cmdParams[0]), false, stop)

	for {
		item, more = <-res
		if !more {
			return nil
		}
		if err = client.RemoveObject(config.bucket, item.Key); err != nil {
			stop <- struct{}{}
			return fmt.Errorf("emptying directory [%s]\nfailed removing the file %s\n%v", config.cmdParams[0], item.Key, err)
		}
	}
	// return nil
}

// downloads a file from the remote location
// cli: `binary` `get` `pwd` `remote file` `local file` `bucketName` `username`
// passed to this is ["remote", "local"]
func get(client minio.CloudStorageClient, config params) error {

	remotePath := expandPath(config.pwd, config.cmdParams[0])

	//func (c Client) FGetObject(bucketName, objectName, filePath string) error
	err := client.FGetObject(config.bucket, remotePath, config.cmdParams[1])
	if err != nil {
		return fmt.Errorf("failed to get object [%q] to file [%q] - %v", err)
	}
	return nil
}

// puts a file onto the remote location
// cli: `binary` `get` `pwd` `local file` `remote file` `bucketName` `username`
// passed to this is ["local", "remote"]
func put(client minio.CloudStorageClient, config params) error {
	contentType, err := fileContentType(config.cmdParams[0])
	if err != nil {
		return fmt.Errorf("failed to determine content type - %v", err)
	}

	remotePath := expandPath(config.pwd, config.cmdParams[1])

	// func (c Client) FPutObject(bucketName, objectName, filePath, contentType string) (n int64, err error)
	_, err = client.FPutObject(config.bucket, remotePath, config.cmdParams[0], contentType)
	if err != nil {
		return fmt.Errorf("failed to put file - %v", err)
	}
	return nil
}

// makes a folder on the remote server.
// cli: `binary` `get` `pwd` `folder name` `bucketName` `username`
// sitting in config.cmdParams is ["folder name"]
func mkdir(client minio.CloudStorageClient, config params) error {
	contents := strings.NewReader("")

	folderName := expandPath(config.pwd, config.cmdParams[0])

	_, err := client.PutObject(config.bucket, folderName+relPathSeperator, contents, folderContentType)
	if err != nil {
		return fmt.Errorf("failed to create folder [%q] - %v", folderName, err)
	}
	return nil
}

func main() {
	config := getConfig(os.Args)

	client, err := getClient(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	switch config.command {
	case "ls":
		err = lsdir(client, config, os.Stdout)
	case "mkdir":
		err = mkdir(client, config)
	case "chdir":
		err = chdir(client, config, os.Stdout)
	case "rmdir":
		err = rmdir(client, config)
	case "delete":
		err = delete(client, config)
	case "get":
		err = get(client, config)
	case "put":
		err = put(client, config)
	default:
		log.Fatal("bad action")
	}

	if err != nil {
		log.Fatal(err.Error())
	}
}
