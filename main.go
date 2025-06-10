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

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// your custom funcs
	fm := template.FuncMap{
		"isTest": isTest,
	}

	// find your layouts
	publicLayouts, err := filepath.Glob(filepath.Join(templatesDir, "layouts", "public.html"))
	if err != nil {
		panic(err)
	}
	adminLayouts, err := filepath.Glob(filepath.Join(templatesDir, "layouts", "admin.html"))
	if err != nil {
		panic(err)
	}

	// walk every file under templatesDir
	err = filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// skip your layout files themselves
		relToLayouts := filepath.Join("layouts")
		if strings.HasPrefix(path, filepath.Join(templatesDir, relToLayouts)) {
			return nil
		}

		// only handle html files
		if !strings.HasSuffix(path, ".html") {
			return nil
		}

		// compute the name you want to render by stripping the base dir
		relPath, err := filepath.Rel(templatesDir, path)
		if err != nil {
			return err
		}
		// turn Windows \ into forward slash
		name := filepath.ToSlash(relPath)

		// pick which layout slice applies
		var layoutSlice []string
		if strings.HasPrefix(relPath, "admin/") {
			layoutSlice = adminLayouts
		} else if strings.HasPrefix(relPath, "public/") {
			layoutSlice = publicLayouts
		}

		// finally register it
		files := append(layoutSlice, path)
		r.AddFromFilesFuncs(name, fm, files...)
		return nil
	})
	if err != nil {
		panic(err)
	}

	return r
}

func setupGin() {
    // Create a new Gin router
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.HTMLRender = loadTemplates("./templates")

	setupRoutes(router)

	// Run the server
	router.Run(":8080")
}

func main() {
	initDB()

	seed()

	setupGin()
}
