package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	. "github.com/yah01/CyDrive/consts"
	"github.com/yah01/CyDrive/model"
	"github.com/yah01/CyDrive/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func LoginHandle(c *gin.Context) {
	username, ok := c.GetPostForm("username")
	if !ok {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusAuthError,
			Message: "no user name",
			Data:    nil,
		})
		return
	}

	password, ok := c.GetPostForm("password")
	if !ok {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusAuthError,
			Message: "no password",
			Data:    nil,
		})
		return
	}

	user := userStore.GetUserByName(username)
	if user == nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusAuthError,
			Message: "no such user",
			Data:    nil,
		})
		return
	}
	if utils.PasswordHash(user.Password) != password {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusAuthError,
			Message: "user name or password not correct",
			Data:    nil,
		})
		return
	}

	userSession := sessions.DefaultMany(c, "user")

	userSession.Set("userStruct", &user)
	userSession.Set("expire", time.Now().Add(time.Hour*12))
	err := userSession.Save()
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: "Welcome to CyDrive!",
		Data:    nil,
	})
}

func ListHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	path := c.Query("path")
	path = strings.Trim(path, string(os.PathSeparator))
	absPath := filepath.Join(user.RootDir, path)

	fileList, err := ioutil.ReadDir(absPath)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	resp := make([]model.FileInfo, 0, len(fileList))
	for _, file := range fileList {
		resp = append(resp, model.NewFileInfo(file,
			filepath.Join(path, file.Name())))
	}
	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: "list done",
		Data:    resp,
	})
}

func GetFileInfoHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	filePath := c.Query("filepath")
	filePath = strings.Trim(filePath, string(os.PathSeparator))
	absFilePath := filepath.Join(user.RootDir, filePath)

	fileInfo, err := os.Stat(absFilePath)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: "get file info done",
		Data:    model.NewFileInfo(fileInfo, filePath),
	})
}

func DownloadHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	// relative path
	filePath := c.Query("filepath")

	// absolute filepath
	filePath = filepath.Join(user.RootDir, filePath)
	fileinfo, _ := os.Stat(filePath)
	if fileinfo.IsDir() {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: "not a file",
			Data:    nil,
		})
		return
	}

	// range
	var begin, end int64 = 0, fileinfo.Size() - 1
	bytesRange := c.GetHeader("Range")
	if len(bytesRange) > 0 {
		begin, end = utils.UnpackRange(bytesRange)
	}

	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	defer file.Close()

	if _, err = file.Seek(begin, io.SeekStart); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.Header("Range", utils.PackRange(begin, end))
	if _, err := io.CopyN(c.Writer, file, end-begin+1); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
}

func UploadHandle(c *gin.Context) {
	if c.Request.ContentLength > FileSizeLimit {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusFileTooLargeError,
			Message: "file is too large",
			Data:    nil,
		})
		return
	}

	userI, _ := c.Get("user")
	user := userI.(*model.User)

	fileInfoJson, ok := c.GetQuery("fileinfo")
	if !ok {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusNoParameterError,
			Message: "need file info",
			Data:    nil,
		})
		return
	}

	fileInfo := model.FileInfo{}
	if err := json.Unmarshal([]byte(fileInfoJson), &fileInfo); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: "error when parsing file info",
			Data:    nil,
		})
		return
	}

	filePath := filepath.Join(user.RootDir, fileInfo.FilePath)
	fileDir := filepath.Dir(filePath)
	if err := os.MkdirAll(fileDir, 0777); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	saveFile, err := os.OpenFile(filePath,
		os.O_RDWR|os.O_CREATE, os.FileMode(fileInfo.FileMode))
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if n, err := io.Copy(saveFile, c.Request.Body); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: fmt.Sprintf("written %v bytes,err: %s", n, err),
			Data:    nil,
		})
		return
	}

	if err = saveFile.Chmod(os.FileMode(fileInfo.FileMode)); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	saveFile.Close()

	if err = os.Chtimes(filePath, time.Now(), time.Unix(fileInfo.ModifyTime, 0)); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: "upload done",
		Data:    nil,
	})
}

func ChangeDirHandle(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(*model.User)

	path := c.Query("path")
	path = strings.Trim(path, string(os.PathSeparator))
	mkdir := c.Query("mkdir")

	var err error

	path = filepath.Join(user.RootDir, path)
	if mkdir == "1" {
		if err = os.MkdirAll(path, 0666); err != nil {
			c.JSON(http.StatusOK, model.Resp{
				Status:  StatusInternalError,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
	}

	_, err = os.Stat(path)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	user.WorkDir = strings.TrimPrefix(path, user.RootDir+"/")

	userSession := sessions.DefaultMany(c, "user")
	userSession.Set("user", user)
	if err = userSession.Save(); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusSessionError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: fmt.Sprintf("you are now in home/%s", user.WorkDir),
		Data:    nil,
	})
}

// The client sends a list consist of all files containing modification time and md5
//func SyncHandle(c *gin.Context) {
//	bodyScanner := bufio.NewScanner(c.Request.Body)
//
//	for bodyScanner.Scan() {
//		line := bodyScanner.Text()
//
//		splitStr := strings.Split(line, " ")
//		filemode, _ := strconv.ParseUint(splitStr[0], 10, 32)
//		modtime, _ := time.Parse(utils.TimeFormat, splitStr[1])
//		filename := splitStr[2]
//
//		os.Create(filename)
//	}
//}
