package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	. "github.com/yah01/CyDrive/consts"
	"github.com/yah01/CyDrive/model"
	"github.com/yah01/CyDrive/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func LoginHandle(c *gin.Context) {
	username, ok := c.GetPostForm("username")
	if !ok {
		c.String(StatusAuthError, "no user name")
		return
	}

	password, ok := c.GetPostForm("password")
	if !ok {
		c.String(StatusAuthError, "no user name")
		return
	}

	user := userStore.GetUserByName(username)
	if user == nil {
		c.String(StatusAuthError, "no such user")
		return
	}
	if utils.PasswordHash(user.Password) != password {
		c.String(StatusAuthError, "user name or password not correct")
		return
	}

	userSession := sessions.DefaultMany(c, "user")

	userSession.Set("userStruct", &user)
	userSession.Set("expire", time.Now().Add(time.Hour*12))
	err := userSession.Save()
	if err != nil {
		c.String(StatusInternalError, "%s", err)
		return
	}
	c.String(StatusOk, "Welcome to CyDrive!")
}

func ListHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	path := c.Query("path")
	path = strings.Trim(path, string(os.PathSeparator))
	path = filepath.Join(user.RootDir, path)

	fileList, err := ioutil.ReadDir(path)
	if err != nil {
		c.String(StatusIoError, "%s", err)
		return
	}

	for _, file := range fileList {
		c.String(StatusOk,
			fmt.Sprintln(file.Mode(), file.ModTime(), file.Name()))
	}
}

func DownloadHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	// relative path
	filePath := c.Query("filepath")

	// absolute filepath
	filePath = filepath.Join(user.RootDir, filePath)
	fileinfo, _ := os.Stat(filePath)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.String(StatusIoError, "%s", err)
		return
	}

	c.JSON(StatusOk, model.File{
		FileInfo: model.FileInfo{
			FileMode:   uint32(fileinfo.Mode()),
			ModifyTime: fileinfo.ModTime().Unix(),
			FilePath:   filePath,
		},
		Data: data,
	})
}

func UploadHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	filePath := c.Query("filepath")
	filePath = filepath.Join(user.RootDir, filePath)

	saveFile, err := os.Create(filePath)
	if err != nil {
		c.String(StatusIoError, "%s", err)
		return
	}

	file := model.File{}
	decoder := json.NewDecoder(c.Request.Body)
	if err = decoder.Decode(&file); err != nil {
		c.String(StatusInternalError, "%s", err)
		return
	}

	if _, err := saveFile.Write(file.Data); err != nil {
		c.String(StatusInternalError, "%s", err)
		return
	}

	if err = saveFile.Chmod(os.FileMode(file.FileMode)); err != nil {
		c.String(StatusInternalError, "%s", err)
		return
	}
	if err = saveFile.Close(); err != nil {
		c.String(StatusInternalError, "%s", err)
		return
	}

	if err = os.Chtimes(filePath, time.Now(), time.Unix(file.ModifyTime, 0)); err != nil {
		c.String(StatusInternalError, "%s", err)
		return
	}

	c.String(StatusOk, "upload %s done", file.FilePath)
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
			c.String(StatusInternalError, "%s", err)
			return
		}
	}

	_, err = os.Stat(path)
	if err != nil {
		c.String(StatusInternalError, "%s", err)
		return
	}

	user.WorkDir = strings.TrimPrefix(path, user.RootDir+"/")

	userSession := sessions.DefaultMany(c, "user")
	userSession.Set("user", user)
	if err = userSession.Save(); err != nil {
		c.String(StatusInternalError, "%s", err)
		return
	}

	c.String(StatusOk, "Done")
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