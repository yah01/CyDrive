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

	server.router.GET("/list", server.ListHandle)

	server.router.GET("/file_info", server.GetFileInfoHandle)
	server.router.PUT("/file_info", server.PutFileInfoHandle)

	server.router.GET("/file", server.GetFileHandle)
	server.router.PUT("/file", server.PutFileHandle)
	server.router.DELETE("/file", server.DeleteFileHandle)

	server.router.Run(ListenPort)
}
