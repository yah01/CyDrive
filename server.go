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

type Server struct {
	Config    config.Config
	userStore store.UserStore
	router    *gin.Engine
}

func NewServer(config config.Config) *Server {
	server := Server{
		Config: config,
	}

	if config.UserStoreType == "mem" {
		server.userStore = store.NewMemStore("user_data/user.json")
	}

	server.router = gin.Default()
	gob.Register(&model.User{})
	gob.Register(time.Time{})

	return &server
}

func (server *Server) Run() {
	memStore := memstore.NewStore([]byte("ProjectMili"))

	server.router.Use(sessions.SessionsMany([]string{"user"}, memStore))
	server.router.Use(LoginAuth(server.router))
	//bin.router.Use(SetFileInfo())

	server.router.POST("/login", server.LoginHandle)

	listGroup := server.router.Group("/list")
	listGroup.GET("/*path", server.ListHandle)

	fileInfoGroup := server.router.Group("/file_info")
	fileInfoGroup.GET("/*path",server.GetFileInfoHandle)
	fileInfoGroup.PUT("/*path", server.PutFileInfoHandle)

	fileGroup := server.router.Group("/file")
	fileGroup.GET("/*path", server.GetFileHandle)
	fileGroup.PUT("/*path", server.PutFileHandle)
	fileGroup.DELETE("/*path", server.DeleteFileHandle)

	server.router.Run(ListenPort)
}
