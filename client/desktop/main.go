package main

import (
	"fmt"
	"github.com/yah01/CyDrive/client"
	. "github.com/yah01/CyDrive/client/consts"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	cli  *client.Client
)

func init() {
	BasePath, _ = exec.LookPath(os.Args[0])
	BasePath, _ = filepath.Abs(BasePath)
	BasePath = filepath.Dir(BasePath)

	BaseUrl = fmt.Sprintf("http://%s:6454", "127.0.0.1")

	cli = client.NewClient()
}

func main() {
	if !cli.Login() {
		time.Sleep(time.Second)
		cli.Login()
	}
	cli.ListRemoteDir("")
	cli.DownloadSync()
	cli.UploadSync()
}
