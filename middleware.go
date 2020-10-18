package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	. "github.com/yah01/CyDrive/consts"
	"github.com/yah01/CyDrive/model"
	"net/http"
	"strings"
	"time"
)

func LoginAuth(router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Trim(c.Request.URL.Path, "/") == "login" {
			c.Next()
			return
		}

		userSession := sessions.DefaultMany(c, "user")
		user := userSession.Get("userStruct")
		expire := userSession.Get("expire")
		if user == nil || expire == nil {
			c.AbortWithStatusJSON(http.StatusOK, model.PackResp(
				StatusAuthError,
				"not login",
				nil,
			))
			return
		}

		if !expire.(time.Time).After(time.Now()) {
			c.AbortWithStatusJSON(http.StatusOK, model.PackResp(
				StatusAuthError,
				"timeout, login again",
				nil,
			))
			userSession.Clear()
			return
		}

		// Flush expire time
		userSession.Set("expire", time.Now().Add(time.Hour*12))
		if err := userSession.Save(); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, model.PackResp(
				StatusSessionError,
				err.Error(),
				nil,
			))
			return
		}

		// Store user struct into context
		c.Set("user", user)
	}
}

func SetFileInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileinfoStr := c.Query("fileinfo")
		if len(fileinfoStr) > 0 {
			fileinfo := model.FileInfo{}
			err := json.Unmarshal([]byte(fileinfoStr), &fileinfo)
			if err != nil {
				fmt.Println(err)
			}

			c.Set("fileinfo", fileinfo)
		}
	}
}
