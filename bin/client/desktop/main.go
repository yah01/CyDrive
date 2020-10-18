package main

import (
	"fmt"
	"github.com/yah01/CyDrive/bin/client"
	. "github.com/yah01/CyDrive/bin/client/consts"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	cli *client.Client
)

func init() {
	BasePath, _ = exec.LookPath(os.Args[0])
	BasePath, _ = filepath.Abs(BasePath)
	BasePath = filepath.Dir(BasePath)

	BaseUrl = fmt.Sprintf("http://%s:6454", "123.57.39.79")

	cli = client.NewClient()
}

func main() {
	if !cli.Login() {
		time.Sleep(time.Second)
		cli.Login()
	}
	cli.ListRemoteDir("")

	for {
		cli.DownloadSync()
		cli.UploadSync()
		time.Sleep(time.Hour * 2)
	}
}
