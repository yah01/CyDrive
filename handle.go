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
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (server *Server) LoginHandle(c *gin.Context) {
	username, ok := c.GetPostForm("username")
	if !ok {
		c.JSON(http.StatusOK, model.PackResp(
			StatusAuthError,
			"no user name",
			nil,
		))
		return
	}

	password, ok := c.GetPostForm("password")
	if !ok {
		c.JSON(http.StatusOK, model.PackResp(
			StatusAuthError,
			"no password",
			nil,
		))
		return
	}

	user := server.userStore.GetUserByName(username)
	if user == nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusAuthError,
			"no such user",
			nil,
		))
		return
	}
	if utils.PasswordHash(user.Password) != password {
		c.JSON(http.StatusOK, model.PackResp(
			StatusAuthError,
			"user name or password not correct",
			nil,
		))
		return
	}

	userSession := sessions.DefaultMany(c, "user")

	userSession.Set("userStruct", &user)
	userSession.Set("expire", time.Now().Add(time.Hour*12))
	err := userSession.Save()
	if err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusInternalError,
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.PackResp(
		StatusOk,
		"Welcome to CyDrive!",
		nil,
	))
}

func (server *Server) ListHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	path, _ := url.QueryUnescape(c.Query("path"))
	path = strings.Trim(path, string(os.PathSeparator))
	absPath := filepath.Join(user.RootDir, path)

	fileList, err := ioutil.ReadDir(absPath)
	if err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusIoError,
			err.Error(),
			nil,
		))
		return
	}

	resp := make([]model.FileInfo, 0, len(fileList))
	for _, file := range fileList {
		resp = append(resp, model.NewFileInfo(file,
			filepath.Join(path, file.Name())))
	}

	c.JSON(http.StatusOK, model.PackResp(
		StatusOk,
		"list done",
		resp,
	))
}

func (server *Server) GetFileInfoHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	filePath, _ := url.QueryUnescape(c.Query("path"))
	filePath = strings.Trim(filePath, string(os.PathSeparator))
	absFilePath := filepath.Join(user.RootDir, filePath)

	fileInfo, err := os.Stat(absFilePath)
	if err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusIoError,
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.PackResp(
		StatusOk,
		"get file info done",
		model.NewFileInfo(fileInfo, filePath),
	))
}

func (server *Server) PutFileInfoHandle(c *gin.Context) {

}

func (server *Server) GetFileHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	// relative path
	filePath, _ := url.QueryUnescape(c.Query("path"))

	// absolute filepath
	filePath = filepath.Join(user.RootDir, filePath)
	fileinfo, _ := os.Stat(filePath)
	if fileinfo.IsDir() {
		c.JSON(http.StatusOK, model.PackResp(
			StatusIoError,
			"not a file",
			nil,
		))
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
		c.JSON(http.StatusOK, model.PackResp(
			StatusIoError,
			err.Error(),
			nil,
		))
		return
	}
	defer file.Close()

	if _, err = file.Seek(begin, io.SeekStart); err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusIoError,
			err.Error(),
			nil,
		))
		return
	}

	c.Header("Range", utils.PackRange(begin, end))
	if _, err := io.CopyN(c.Writer, file, end-begin+1); err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusIoError,
			err.Error(),
			nil,
		))
		return
	}
}

func (server *Server) PutFileHandle(c *gin.Context) {
	// Check file size
	if c.Request.ContentLength > FileSizeLimit {
		c.JSON(http.StatusOK, model.PackResp(
			StatusFileTooLargeError,
			"file is too large",
			nil,
		))
		return
	}

	userI, _ := c.Get("user")
	user := userI.(*model.User)

	// Check user storage capability
	if c.Request.ContentLength+user.Usage > user.Cap {
		c.JSON(http.StatusOK, model.PackResp(
			StatusFileTooLargeError,
			fmt.Sprintf("no enough capability, free storage: %vMB",
				(user.Cap-user.Usage)>>20), // Convert Byte to MB
			nil,
		))
		return
	}

	fileInfoJson, ok := c.GetQuery("fileinfo")
	if !ok {
		c.JSON(http.StatusOK, model.PackResp(
			StatusNoParameterError,
			"need file info",
			nil,
		))
		return
	}
	fileInfoJson, _ = url.QueryUnescape(fileInfoJson)

	fileInfo := model.FileInfo{}
	if err := json.Unmarshal([]byte(fileInfoJson), &fileInfo); err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusInternalError,
			"error when parsing file info",
			nil,
		))
		return
	}

	filePath := filepath.Join(user.RootDir, fileInfo.FilePath)
	fileDir := filepath.Dir(filePath)
	if err := os.MkdirAll(fileDir, 0777); err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusInternalError,
			err.Error(),
			nil,
		))
		return
	}

	saveFile, err := os.OpenFile(filePath,
		os.O_RDWR|os.O_CREATE, os.FileMode(fileInfo.FileMode))
	if err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusIoError,
			err.Error(),
			nil,
		))
		return
	}

	if n, err := io.Copy(saveFile, c.Request.Body); err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusIoError,
			fmt.Sprintf("written %v bytes,err: %s", n, err),
			nil,
		))
		return
	}

	if err = saveFile.Chmod(os.FileMode(fileInfo.FileMode)); err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusInternalError,
			err.Error(),
			nil,
		))
		return
	}

	saveFile.Close()

	if err = os.Chtimes(filePath, time.Now(), time.Unix(fileInfo.ModifyTime, 0)); err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusInternalError,
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.PackResp(
		StatusOk,
		"upload done",
		nil,
	))
}

func (server *Server) DeleteFileHandle(c *gin.Context) {

}

func (server *Server) ChangeDirHandle(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(*model.User)

	path, _ := url.QueryUnescape(c.Query("path"))
	path = strings.Trim(path, string(os.PathSeparator))
	mkdir := c.Query("mkdir")

	var err error

	path = filepath.Join(user.RootDir, path)
	if mkdir == "1" {
		if err = os.MkdirAll(path, 0666); err != nil {
			c.JSON(http.StatusOK, model.PackResp(
				StatusInternalError,
				err.Error(),
				nil,
			))
			return
		}
	}

	_, err = os.Stat(path)
	if err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusInternalError,
			err.Error(),
			nil,
		))
		return
	}

	user.WorkDir = strings.TrimPrefix(path, user.RootDir+"/")

	userSession := sessions.DefaultMany(c, "user")
	userSession.Set("user", user)
	if err = userSession.Save(); err != nil {
		c.JSON(http.StatusOK, model.PackResp(
			StatusSessionError,
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.PackResp(
		StatusOk,
		fmt.Sprintf("you are now in home/%s", user.WorkDir),
		nil,
	))
}

// The client sends a list consist of all files containing modification time and md5
//func (server *Server) SyncHandle(c *gin.Context) {
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
