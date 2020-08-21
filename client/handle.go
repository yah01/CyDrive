package main

import (
	"bufio"
	"fmt"
	"github.com/mholt/archiver/v3"
	"github.com/yah01/CyDrive/utils"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func Login(username string, password string) {
	Url, _ := url.Parse(baseUrl + "/login")

	originTimeout := client.Timeout
	client.Timeout = time.Second * 120
	resp, err := client.PostForm(Url.String(), url.Values{
		"username": {username},
		"password": {utils.PasswordHash(password)},
	})
	client.Timeout = originTimeout
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(GetRespInfo(resp))
}

func ListRemoteDir(path ...string) {
	if len(path) == 0 {
		Url, _ := url.Parse(baseUrl + "/list")
		resp, err := client.Get(Url.String())
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		fmt.Println(GetRespInfo(resp))
	}
}

func Download(filename string) {
	Url, _ := url.Parse(fmt.Sprint(baseUrl+"/download?filepath=", filename))
	resp, err := client.Get(Url.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("downloading...")
	if resp.Header.Get("Content-Type") == "dir" {
		zipFileName := "tmp/" + filename + ".zip"
		file, _ := os.Create(zipFileName)
		io.Copy(file, resp.Body)
		file.Close()

		archiver.Unarchive(zipFileName, user.WorkDir)
		os.Remove(zipFileName)
	} else {
		file, _ := os.Create(filepath.Join(user.WorkDir, filename))
		defer file.Close()
		io.Copy(file, resp.Body)
	}
	fmt.Println("done")
}

func Upload(filename string) {
	Url, _ := url.Parse(fmt.Sprint(baseUrl+"/upload?filename=", filename))

	path := filepath.Join(user.WorkDir, filename)
	contentType := "file"
	fileinfo, _ := os.Stat(path)
	if fileinfo.IsDir() {
		zipFileName := "tmp/" + fileinfo.Name() + ".zip"
		archiver.Archive([]string{path}, zipFileName)
		path = zipFileName
		contentType = "dir"
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Println("uploading...")

	resp, err := client.Post(Url.String(), contentType, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("done")
	if fileinfo.IsDir() {
		os.Remove(path)
	}
}

func ChangeRemoteDir(path string) {
	Url, _ := url.Parse(fmt.Sprint(baseUrl+"/change_dir?path=", path))
	resp, err := client.Get(Url.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(GetRespInfo(resp))
}

func GetRespInfo(resp *http.Response) string {
	reader := bufio.NewReader(resp.Body)
	var (
		info, tmp string
		err       error
	)

	for {
		tmp, err = reader.ReadString('\n')
		info = info + tmp
		if err != nil {
			break
		}
	}
	return info
}
