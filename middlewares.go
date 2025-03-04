package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func middlewareAuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("currentUser")

	if userId == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	} else {
		var user User
		if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
			session.Delete("currentUser")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
	}

	c.Next()
}

func middlewareSetUser(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("currentUser")

	if userId == nil {
		c.Next()
		return
	}

	var user User
	if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
		return
	}
	c.Set("currentUser", user)

	c.Next()
}
