package main

import "github.com/gin-gonic/gin"

func setupRoutes(router *gin.Engine) {
	router.GET("/", actionPublicRoot)

	router.GET("/pages/:slug", actionPublicPage)

	router.GET("/login", middlewareSetUser, actionPublicLoginForm)
	router.POST("/login", middlewareSetUser, actionPublicLoginSubmit)
	router.GET("/logout", middlewareSetUser, actionPublicLogout)

	router.GET("/admin", actionAdminIndex)
	router.GET("/admin/users/new", middlewareAuthRequired, middlewareSetUser, actionAdminUsersNew)
	router.POST("/admin/users/create", middlewareAuthRequired, middlewareSetUser, actionAdminUsersCreate)
	router.GET("/admin/users", middlewareAuthRequired, middlewareSetUser, actionAdminUsersIndex)
	router.GET("/admin/users/:id", middlewareAuthRequired, middlewareSetUser, actionAdminUsersShow)
	router.GET("/admin/users/:id/edit", middlewareAuthRequired, middlewareSetUser, actionAdminUsersEdit)
	router.POST("/admin/users/:id/update", middlewareAuthRequired, middlewareSetUser, actionAdminUsersUpdate)
	router.POST("/admin/users/:id/delete", middlewareAuthRequired, middlewareSetUser, actionAdminUsersDestroy)
	router.GET("/admin/pages", middlewareAuthRequired, middlewareSetUser, actionAdminPagesIndex)
	router.GET("/admin/pages/new", middlewareAuthRequired, middlewareSetUser, actionAdminPagesNew)
	router.POST("/admin/pages/create", middlewareAuthRequired, middlewareSetUser, actionAdminPagesCreate)
	router.GET("/admin/pages/:id", middlewareAuthRequired, middlewareSetUser, actionAdminPagesShow)
	router.GET("/admin/pages/:id/edit", middlewareAuthRequired, middlewareSetUser, actionAdminPagesEdit)
	router.POST("/admin/pages/:id/update", middlewareAuthRequired, middlewareSetUser, actionAdminPagesUpdate)
	router.POST("/admin/pages/:id/delete", middlewareAuthRequired, middlewareSetUser, actionAdminPagesDestroy)

	if isTest() {
		router.GET("/tools", actionPublicTools)
		router.GET("/tools/db-clear", actionPublicToolsDBClear)
		router.GET("/tools/seed", actionPublicToolsSeed)
		router.GET("/tools/sql", actionPublicToolsSQL)
	}
}
