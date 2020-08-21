package main

// a simple cydrive client only for test

import (
	"bufio"
	"fmt"
	. "github.com/yah01/CyDrive/const"
	"github.com/yah01/CyDrive/model"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
)

var (
	cookieJar *cookiejar.Jar
	client    *http.Client
	user      = model.User{
		Username: "test",
		Password: "testCyDrive",
	}

	baseUrl string
)

func init() {
	cookieJar, _ = cookiejar.New(nil)
	client = &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           cookieJar,
		Timeout:       0,
	}
}
var serverAddress = "127.0.0.1"

func main() {

	baseUrl = fmt.Sprintf("http://%s:6454", serverAddress)
	Login(user.Username, user.Password)

	var (
		cmd    string
		reader = bufio.NewReader(os.Stdin)
	)

	for {
		cmd, _ = reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		cmdSplit := strings.Split(cmd, " ")
		cmd = strings.ToUpper(cmdSplit[0])

		switch cmd {
		// communicate with server:
		case LOGIN:
			Login(cmdSplit[1], cmdSplit[2])
		case LIST:
			ListRemoteDir()
		case GET:
			Download(cmdSplit[1])
		case SEND:
			Upload(cmdSplit[1])
		case RCD:
			ChangeRemoteDir(cmdSplit[1])
		case QUIT:
			client.CloseIdleConnections()
			return

			// not communicate with server
		}
	}
}