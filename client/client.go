package client

import (
	"encoding/json"
	"fmt"
	. "github.com/yah01/CyDrive/client/consts"
	"github.com/yah01/CyDrive/consts"
	"github.com/yah01/CyDrive/model"
	"github.com/yah01/CyDrive/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	cookieJar  *cookiejar.Jar
	user       model.User
}

func NewClient() *Client {
	client := Client{}
	client.cookieJar, _ = cookiejar.New(nil)
	client.httpClient = &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           client.cookieJar,
		Timeout:       0,
	}
	userBytes, _ := ioutil.ReadFile(filepath.Join(BasePath, "user.json"))
	json.Unmarshal(userBytes, &client.user)
	return &client
}

func (client *Client) Login() bool {
	Url, _ := url.Parse(BaseUrl + "/login")

	originTimeout := client.httpClient.Timeout
	client.httpClient.Timeout = time.Second * 120
	resp, err := client.httpClient.PostForm(Url.String(), url.Values{
		"username": {client.user.Username},
		"password": {utils.PasswordHash(client.user.Password)},
	})
	if err != nil {
		fmt.Println(err)
		return false
	}
	client.httpClient.Timeout = originTimeout
	defer resp.Body.Close()

	res := utils.GetResp(resp)
	fmt.Println(res)
	return res.Status == consts.StatusOk
}

func (client *Client) ListRemoteDir(path string) []*model.FileInfo {
	Url, _ := url.Parse(
		fmt.Sprintf(BaseUrl+"/list?path=%s", path))

	resp, err := client.httpClient.Get(Url.String())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	res := model.Resp{}
	if err = json.Unmarshal(data, &res); err != nil || res.Status != consts.StatusOk {
		return nil
	}

	list := res.Data.([]interface{})
	fileInfoList := make([]*model.FileInfo, 0, len(list))
	for _, file := range list {
		fileInfoMap := file.(map[string]interface{})
		fileInfo := model.NewFileInfoFromMap(fileInfoMap)
		fileInfoList = append(fileInfoList, fileInfo)
		fmt.Println(fileInfo)
	}
	return fileInfoList
}

func (client *Client) GetFileInfo(path string) *model.FileInfo {
	Url, _ := url.Parse(
		fmt.Sprintf(BaseUrl+"/get_file_info?filepath=%s", path))

	resp, err := client.httpClient.Get(Url.String())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	res := utils.GetResp(resp)
	if res.Status != consts.StatusOk {
		return nil
	}
	fileInfo := model.NewFileInfoFromMap(res.Data.(map[string]interface{}))
	return fileInfo
}

func (client *Client) Download(path string) bool {
	fileInfo := client.GetFileInfo(path)

	Url, _ := url.Parse(
		fmt.Sprintf(BaseUrl+"/download?filepath=%s", path))

	resp, err := client.httpClient.Get(Url.String())
	if err != nil {
		fmt.Println(err)
		return false
	}

	absPath := filepath.Join(client.user.RootDir, path)
	err = os.MkdirAll(filepath.Dir(absPath), 0777)
	if err != nil {
		fmt.Println(err)
	}

	saveFile, err := os.OpenFile(absPath,
		os.O_RDWR|os.O_CREATE, os.FileMode(fileInfo.FileMode))
	if err != nil {
		fmt.Println(err)
		return false
	}

	if _, err := io.Copy(saveFile, resp.Body); err != nil {
		fmt.Println(err)
		return false
	}

	if err = saveFile.Chmod(os.FileMode(fileInfo.FileMode)); err != nil {
		fmt.Println(err)
		return false
	}

	saveFile.Close()

	if err = os.Chtimes(absPath, time.Now(), time.Unix(fileInfo.ModifyTime, 0)); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (client *Client) Upload(path string) *model.Resp {
	osFileInfo, _ := os.Stat(filepath.Join(client.user.RootDir, path))
	fileInfo := model.NewFileInfo(osFileInfo, path)
	bytes, _ := json.Marshal(fileInfo)

	Url, _ := url.Parse(
		fmt.Sprintf(BaseUrl+"/upload?fileinfo=%s", string(bytes)),
	)

	localPath := filepath.Join(client.user.RootDir, path)
	file, _ := os.Open(localPath)
	resp, _ := client.httpClient.Post(Url.String(), "file", file)
	defer resp.Body.Close()

	res := utils.GetResp(resp)
	fmt.Println(res)
	return res
}

func (client *Client) DownloadSync() {
	utils.ForEachRemoteFile("",
		client.GetFileInfo, client.ListRemoteDir, func(file *model.FileInfo) {
			localPath := filepath.Join(client.user.RootDir, file.FilePath)
			localFile, err := os.Open(localPath)
			var fileInfo os.FileInfo = nil
			if localFile != nil {
				fileInfo, _ = localFile.Stat()
			}

			if err != nil || localFile == nil || fileInfo.ModTime().Unix() < file.ModifyTime {
				fmt.Println("download sync:", file.FilePath)
				ok := client.Download(file.FilePath)
				if !ok {
					fmt.Println("can't download:", file.FilePath)
				}
			}
		})
}

func (client *Client) UploadSync() {
	utils.ForEachFile(client.user.RootDir, func(file *os.File) {
		remotePath := strings.TrimPrefix(file.Name(), client.user.RootDir+string(os.PathSeparator))
		remoteFile := client.GetFileInfo(remotePath)

		fileInfo, _ := file.Stat()
		if remoteFile == nil || remoteFile.ModifyTime < fileInfo.ModTime().Unix() {
			fmt.Println("upload sync:", file.Name())
			resp := client.Upload(remotePath)
			if resp.Status != 0 {
				fmt.Println(resp)
			}
		}
	})
}
