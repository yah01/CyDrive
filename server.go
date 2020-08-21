package main

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yah01/CyDrive/config"
	. "github.com/yah01/CyDrive/consts"
	"github.com/yah01/CyDrive/model"
	"github.com/yah01/CyDrive/store"
	"time"
)

var (
	userStore store.UserStore
	router    *gin.Engine
)

func InitServer(config config.Config) {
	if config.UserStoreType == "mem" {
		userStore = store.NewMemStore("user_data/user.json")
	}

	router = gin.Default()
	gob.Register(&model.User{})
	gob.Register(time.Time{})
}

func RunServer() {
	memStore := memstore.NewStore([]byte("ProjectMili"))

	router.Use(sessions.SessionsMany([]string{"user"}, memStore))
	router.Use(LoginAuth(router))
	//router.Use(SetFileInfo())

	router.POST("/login", LoginHandle)
	router.GET("/list", ListHandle)
	router.GET("/download", DownloadHandle)
	router.POST("/upload", UploadHandle)
	router.GET("/change_dir", ChangeDirHandle)
	//router.GET("/sync",SyncHandle)
	router.Run(ListenPort)
}
