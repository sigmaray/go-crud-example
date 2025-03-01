package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	db.AutoMigrate(&User{})
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

		for _, publicTemplate := range publicTemplates {
			layoutCopy := make([]string, len(publicLayouts))
			copy(layoutCopy, publicLayouts)
			files := append(layoutCopy, publicTemplate)
			r.AddFromFiles(filepath.Base(publicTemplate), files...)
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
			r.AddFromFiles(filepath.Base(adminTemplate), files...)
		}
		return r
	}

	// Create a new Gin router
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.HTMLRender = loadTemplates("./templates")

	router.GET("/", actionRoot)

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

	// Run the server
	router.Run(":8080")
}

func main() {
	initDB()

	seed()

	setupGin()
}
