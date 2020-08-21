package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	. "github.com/yah01/CyDrive/const"
	"github.com/yah01/CyDrive/model"
	"strings"
	"time"
)

func LoginAuth(router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Trim(c.Request.URL.Path, "/") == "login" {
			return
		}

		userSession := sessions.DefaultMany(c, "user")
		user := userSession.Get("userStruct")
		expire := userSession.Get("expire")
		if user == nil || expire == nil {
			c.String(StatusAuthError, "not login")
			return
		}

		if !expire.(time.Time).After(time.Now()) {
			c.String(StatusAuthError, "timeout, login again")
			userSession.Clear()
			return
		}

		// flush expire time
		userSession.Set("expire", time.Now().Add(time.Hour*12))
		userSession.Save()

		// store user struct into context
		c.Set("user", user)
	}
}

func SetFileInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileinfoStr := c.Query("fileinfo")
		if len(fileinfoStr) > 0 {
			fileinfo := model.FileInfo{}
			err := json.Unmarshal([]byte(fileinfoStr),&fileinfo)
			if err != nil {
				fmt.Println(err)
			}

			c.Set("fileinfo",fileinfo)
		}
	}
}