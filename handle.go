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
	path = filepath.Join(user.RootDir, path)

	fileList, err := ioutil.ReadDir(path)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	resp := make([]string, 0, len(fileList))
	for _, file := range fileList {
		resp = append(resp,
			fmt.Sprintf("%s %s %s", file.Mode(), file.ModTime(), file.Name()))
	}
	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: "list done",
		Data:    resp,
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

	data, err := ioutil.ReadFile(filePath)
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
		Message: "download done",
		Data: model.File{
			FileInfo: model.FileInfo{
				FileMode:   uint32(fileinfo.Mode()),
				ModifyTime: fileinfo.ModTime().Unix(),
				FilePath:   filePath,
			},
			Data: data,
		},
	})
}

func UploadHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	filePath := c.Query("filepath")
	filePath = filepath.Join(user.RootDir, filePath)

	saveFile, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	file := model.File{}
	decoder := json.NewDecoder(c.Request.Body)
	if err = decoder.Decode(&file); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if _, err := saveFile.Write(file.Data); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err = saveFile.Chmod(os.FileMode(file.FileMode)); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if err = saveFile.Close(); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err = os.Chtimes(filePath, time.Now(), time.Unix(file.ModifyTime, 0)); err != nil {
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
			Status:  StatusInternalError,
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
