package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver/v3"
	. "github.com/yah01/CyDrive/const"
	"github.com/yah01/CyDrive/model"
	"github.com/yah01/CyDrive/utils"
	"io"
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
	if user==nil {
		c.String(StatusAuthError,"no such user")
		return
	}
	if utils.PasswordHash(user.Password) != password {
		c.String(StatusAuthError,"user name or password not correct")
		return
	}

	userSession := sessions.DefaultMany(c, "user")

	userSession.Set("userStruct", &user)
	userSession.Set("expire", time.Now().Add(time.Hour*12))
	err := userSession.Save()
	if err != nil {
		fmt.Println(err)
	}
	c.String(StatusOk, "Welcome to CyDrive!")
}

func ListHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)
	fileList, _ := ioutil.ReadDir(JoinUserPath(user))

	for _, file := range fileList {
		c.String(StatusOk,
			fmt.Sprintln(file.Mode(), file.ModTime(), file.Name()))
	}
}

func DownloadHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	// Local relative path
	filename := c.Query("filepath")

	// Remote filepath
	filePath := JoinUserPath(user, filepath.Base(filename))
	fileinfo, _ := os.Stat(filePath)
	if fileinfo.IsDir() {
		zipFileName := "tmp/" + fileinfo.Name() + ".zip"

		archiver.Archive([]string{filePath}, zipFileName)
		defer os.Remove(zipFileName)

		c.Header("Content-Type", "dir")
		c.File(zipFileName)
	} else {
		data, _ := ioutil.ReadFile(filePath)
		file := model.File{
			FileInfo: model.FileInfo{
				FileMode:   uint32(fileinfo.Mode()),
				ModifyTime: fileinfo.ModTime().Unix(),
				FilePath:   filename,
			},
			Data: data,
		}

		c.JSON(StatusOk, file)
	}
}

func UploadHandle(c *gin.Context) {
	user, _ := c.Get("user")

	filename := c.Query("filename")
	contentType := c.GetHeader("Content-Type")

	path := JoinUserPath(user.(*model.User), filename)
	if contentType == "dir" {
		zipFileName := "tmp/" + filename + ".zip"
		file, _ := os.Create(zipFileName)
		io.Copy(file, c.Request.Body)
		file.Close()

		archiver.Unarchive(zipFileName, JoinUserPath(user.(*model.User)))
		os.Remove(zipFileName)

	} else {
		file, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		io.Copy(file, c.Request.Body)
	}
	c.String(StatusOk, "saved file %s", filename)
}

func ChangeDirHandle(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(*model.User)

	path := c.Query("path")
	path = strings.Trim(path, string(os.PathSeparator))

	if strings.HasPrefix(path, "~") {
		path = strings.ReplaceAll(path, "~", "")
		path = strings.Trim(path, string(os.PathSeparator))

		absPath := filepath.Join(user.RootDir, path)
		fileinfo, _ := os.Stat(absPath)
		if !fileinfo.IsDir() {
			c.String(StatusOk, "no such dir")
			return
		}
		user.WorkDir = strings.Trim(strings.TrimPrefix(absPath, user.RootDir), string(os.PathSeparator))
		c.String(StatusOk, "now you are in %s", user.WorkDir)
	} else if strings.HasPrefix(path, "..") {
		dotDotCount := 0
		for strings.HasPrefix(path, "..") {
			path = strings.TrimPrefix(path, "..")
			path = strings.Trim(path, string(os.PathSeparator))
			dotDotCount++
		}

		pathList := strings.Split(user.WorkDir, string(os.PathSeparator))
		if dotDotCount > len(pathList) {
			c.String(StatusOk, "wrong path")
			return
		}

		pathList = pathList[:len(pathList)-dotDotCount]
		pathList = append(pathList,
			strings.Split(path, string(os.PathSeparator))...)
		pathList = append([]string{user.RootDir}, pathList...)
		absPath, err := filepath.Abs(filepath.Join(pathList...))
		if err != nil {
			c.String(StatusOk, "wrong path")
			return
		}

		fileinfo, _ := os.Stat(absPath)
		if !fileinfo.IsDir() {
			c.String(StatusOk, "wrong path")
			return
		}

		user.WorkDir = strings.Trim(strings.TrimPrefix(absPath, user.RootDir), string(os.PathSeparator))
		c.String(StatusOk, "now you are in %s", user.WorkDir)
	} else {
		path = strings.TrimPrefix(path, ".")
		path = strings.TrimPrefix(path, string(os.PathSeparator))

		path = JoinUserPath(user, path)
		absPath, err := filepath.Abs(path)

		if err != nil {
			c.String(StatusOk, "wrong path")
			return
		}

		fileinfo, _ := os.Stat(absPath)
		if !fileinfo.IsDir() {
			c.String(StatusOk, "wrong path")
			return
		}

		user.WorkDir = strings.Trim(strings.TrimPrefix(absPath, user.RootDir), string(os.PathSeparator))
		c.String(StatusOk, "now you are in %s", user.WorkDir)
	}

	userSession := sessions.DefaultMany(c, "user")
	userSession.Set("user", user)
	userSession.Save()
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

func JoinUserPath(user *model.User, path ...string) string {
	res := []string{UserDataDir, fmt.Sprint(user.Id), user.WorkDir}
	res = append(res, path...)
	return filepath.Join(res...)
}
