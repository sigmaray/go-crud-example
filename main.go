package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var err error

// isDocker checks if the program is running inside a Docker container
func isDocker() bool {
	// Check if the /.dockerenv file exists
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// Check if the /proc/self/cgroup file exists and contains 'docker'
	cgroupFile := "/proc/self/cgroup"
	if _, err := os.Stat(cgroupFile); err == nil {
		content, err := os.ReadFile(cgroupFile)
		if err == nil && strings.Contains(string(content), "docker") {
			return true
		}
	}

	return false
}

// TODO: Move DB credentials into .env
func initDB() {
	// Connect to the database
	var pg_host string
	if isDocker() {
		pg_host = "postgres"
	} else {
		pg_host = "localhost"
	}
	dsn := "host=" + pg_host + " user=postgres dbname=appdb password=postgres sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	db.AutoMigrate(&User{}, &Page{})
}

func seed() {
	var count int64

	db.Model(&User{}).Count(&count)

	// If there are no users, create sample user
	if count == 0 {
		result := db.Create(&User{Login: "admin", Password: "admin"})
		if result.Error != nil {
			println(result.Error)
			panic("Failed to create user")
		}

	}
}

func isTest() bool {
	return (os.Getenv("TEST") == "1" || os.Getenv("TEST") == "true" || os.Getenv("TEST") == "yes" || os.Getenv("TEST") == "on" || os.Getenv("TEST") == "t")
}

func setupGin() {
	loadTemplates := func(templatesDir string) multitemplate.Renderer {
		r := multitemplate.NewRenderer()

		publicLayouts, err := filepath.Glob(templatesDir + "/layouts/public.html")
		if err != nil {
			panic(err.Error())
		}

		publicTemplates, err := filepath.Glob(templatesDir + "/*.html")
		if err != nil {
			panic(err.Error())
		}

		fm := template.FuncMap{
			"isTest": func() bool {
				return isTest()
			},
		}

		for _, publicTemplate := range publicTemplates {
			layoutCopy := make([]string, len(publicLayouts))
			copy(layoutCopy, publicLayouts)
			files := append(layoutCopy, publicTemplate)
			r.AddFromFilesFuncs(filepath.Base(publicTemplate), fm, files...)
		}

		adminLayouts, err := filepath.Glob(templatesDir + "/layouts/admin.html")
		if err != nil {
			panic(err.Error())
		}

		admins, err := filepath.Glob(templatesDir + "/admin/*.html")
		if err != nil {
			panic(err.Error())
		}

		for _, adminTemplate := range admins {
			layoutCopy := make([]string, len(adminLayouts))
			copy(layoutCopy, adminLayouts)
			files := append(layoutCopy, adminTemplate)
			r.AddFromFilesFuncs(filepath.Base(adminTemplate), fm, files...)
		}
		return r
	}

	// Create a new Gin router
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.HTMLRender = loadTemplates("./templates")

	router.GET("/", actionRoot)

	router.GET("/pages/:slug", actionPage)

	router.GET("/login", middlewareSetUser, actionLoginForm)
	router.POST("/login", middlewareSetUser, actionLoginSubmit)
	router.GET("/logout", middlewareSetUser, actionLogout)

	router.GET("/admin", actionAdminIndex)
	router.GET("/admin/users/new", middlewareAuthRequired, middlewareSetUser, actionUsersNew)
	router.POST("/admin/users/create", middlewareAuthRequired, middlewareSetUser, actionUsersCreate)
	router.GET("/admin/users", middlewareAuthRequired, middlewareSetUser, actionUsersIndex)
	router.GET("/admin/users/:id", middlewareAuthRequired, middlewareSetUser, actionUsersShow)
	router.GET("/admin/users/:id/edit", middlewareAuthRequired, middlewareSetUser, actionUsersEdit)
	router.POST("/admin/users/:id/update", middlewareAuthRequired, middlewareSetUser, actionUsersUpdate)
	router.POST("/admin/users/:id/delete", middlewareAuthRequired, middlewareSetUser, actionUsersDestroy)
	router.GET("/admin/pages", middlewareAuthRequired, middlewareSetUser, actionPagesIndex)
	router.GET("/admin/pages/new", middlewareAuthRequired, middlewareSetUser, actionPagesNew)
	router.POST("/admin/pages/create", middlewareAuthRequired, middlewareSetUser, actionPagesCreate)
	router.GET("/admin/pages/:id", middlewareAuthRequired, middlewareSetUser, actionPagesShow)
	router.GET("/admin/pages/:id/edit", middlewareAuthRequired, middlewareSetUser, actionPagesEdit)
	router.POST("/admin/pages/:id/update", middlewareAuthRequired, middlewareSetUser, actionPagesUpdate)
	router.POST("/admin/pages/:id/delete", middlewareAuthRequired, middlewareSetUser, actionPagesDestroy)
	router.GET("/tools", actionTools)
	router.GET("/tools/db-clear", actionToolsDBClear)
	router.GET("/tools/seed", actionToolsSeed)
	router.GET("/tools/sql", actionToolsSQL)

	// router.GET("/admin/collection_todos/new", middlewareAuthRequired, middlewareSetUser, actionCollectionTodosNew)
	// router.POST("/admin/collection_todos/create", middlewareAuthRequired, middlewareSetUser, actionCollectionTodosCreate)
	router.GET("/admin/collection_todos", middlewareAuthRequired, middlewareSetUser, actionCollectionTodosIndex)

	// Run the server
	router.Run(":8080")
}

func main() {
	initDB()

	seed()

	setupGin()
}
