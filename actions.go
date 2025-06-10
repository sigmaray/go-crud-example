package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Convert validation errors into slice of human readable error strings
func humanValidationErrors(err error) []string {
	getErrorMessage := func(field, tag string) string {
		if tag == "required" {
			return "[Validation error] " + field + ": Field is required\n"
		}

		if tag == "min" {
			return "[Validation error] " + field + ": Field is too short\n"
		}

		return "[Validation error] " + field + ": Invalid input\n"
	}

	var errorMessages []string
	var validateErrs validator.ValidationErrors
	errors.As(err, &validateErrs)

	for _, err := range validateErrs {
		field := err.Field()
		tag := err.Tag()
		errorMessages = append(errorMessages, getErrorMessage(field, tag))
	}
	return errorMessages
}

func addFlashesAndUser(c *gin.Context, h *gin.H) *gin.H {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()
	(*h)["flashes"] = flashes

	user, _ := c.Get("currentUser")
	if user != nil {
		(*h)["currentUser"] = user.(User)
	} else {
		(*h)["currentUser"] = nil
	}

	return h
}

func actionPublicRoot(c *gin.Context) {
	var pages []Page
	db.Find(&pages)
	c.HTML(http.StatusOK, "public/index.html", &gin.H{"pages": pages})
}

func actionPublicPage(c *gin.Context) {
	slug := c.Param("slug")

	var page Page

	if err := db.Where("slug = ?", slug).First(&page).Error; err != nil {
		c.HTML(http.StatusNotFound, "public/error.html", addFlashesAndUser(c, &gin.H{"errors": []string{"Page not found"}}))
		return
	}

	pageJSON, err := json.MarshalIndent(page, "", "  ")
	if err != nil {
		c.HTML(http.StatusNotFound, "public/error.html", addFlashesAndUser(c, &gin.H{"errors": []string{err.Error()}}))
		return
	}

	c.HTML(http.StatusOK, "public/page.html", &gin.H{"slug": slug, "page": page, "pageJSON": string(pageJSON)})
}

func actionPublicLoginForm(c *gin.Context) {
	_, exists := c.Get("currentUser")
	if exists {
		c.Redirect(http.StatusSeeOther, "/admin/users")
	}

	c.HTML(http.StatusOK, "public/login.html", nil)
}

func actionPublicLoginSubmit(c *gin.Context) {
	_, exists := c.Get("currentUser")
	if exists {
		c.Redirect(http.StatusSeeOther, "/admin/users")
	}

	username := c.PostForm("login")
	password := c.PostForm("password")

	var user User
	// TODO: Encrypt password
	if err := db.Where("login = ? and password = ?", username, password).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "public/login.html", gin.H{"errors": []string{"Invalid username or password"}})
		return
	}

	session := sessions.Default(c)
	session.Set("currentUser", user.ID)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/users")
}

func actionPublicLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("currentUser")
	session.AddFlash("Logged out")
	session.Save()
	c.Redirect(http.StatusSeeOther, "/admin/users")
}

func actionAdminUsersIndex(c *gin.Context) {
	var users []User
	db.Find(&users)
	c.HTML(http.StatusOK, "admin/users/index.html", addFlashesAndUser(c, &gin.H{"users": users}))
}

func actionAdminUsersShow(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := db.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "admin/error.html", addFlashesAndUser(c, &gin.H{"errors": []string{"User not found"}}))
		return
	}

	userJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		c.HTML(http.StatusNotFound, "admin/error.html", addFlashesAndUser(c, &gin.H{"errors": []string{err.Error()}}))
		return
	}

	c.HTML(http.StatusOK, "admin/users/show.html", addFlashesAndUser(c, &gin.H{"user": user, "userJSON": string(userJSON)}))
}

func actionAdminIndex(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/admin/users")
}

func actionAdminUsersNew(c *gin.Context) {
	var user User
	c.HTML(http.StatusOK, "admin/users/new.html", addFlashesAndUser(c, &gin.H{"user": user}))
}

func actionAdminUsersCreate(c *gin.Context) {
	var user User
	user.Login = c.PostForm("login")
	// TODO: Encrypt password
	user.Password = c.PostForm("password")
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	user_input := &UserInput{
		Login:    user.Login,
		Password: user.Password,
	}

	// Validate user input
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(user_input); err != nil {
		c.HTML(http.StatusBadRequest, "admin/users/new.html", addFlashesAndUser(c, &gin.H{"errors": humanValidationErrors(err), "user": user}))
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin/users/new.html", addFlashesAndUser(c, &gin.H{"errors": []string{err.Error()}, "user": user}))
		return
	}

	session := sessions.Default(c)
	session.AddFlash("User was added.")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/users")
}

func actionAdminUsersEdit(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := db.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "admin/error.html", addFlashesAndUser(c, &gin.H{"errors": []string{"User not found"}}))
		return
	}

	c.HTML(http.StatusOK, "admin/users/edit.html", addFlashesAndUser(c, &gin.H{"user": user}))
}

func actionAdminUsersUpdate(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := db.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "admin/error.html", addFlashesAndUser(c, &gin.H{"errors": []string{"User not found"}}))
		return
	}

	user.Login = c.PostForm("login")
	user.Password = c.PostForm("password")
	user.UpdatedAt = time.Now()

	user_input := &UserInput{
		Login:    user.Login,
		Password: user.Password,
	}

	// Validate user input
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(user_input); err != nil {
		c.HTML(http.StatusOK, "admin/users/edit.html", addFlashesAndUser(c, &gin.H{"errors": humanValidationErrors(err), "user": user}))
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.HTML(http.StatusOK, "admin/users/edit.html", addFlashesAndUser(c, &gin.H{"errors": []string{err.Error()}, "user": user}))
		return
	}

	session := sessions.Default(c)
	session.AddFlash("User was edited.")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/users")
}

func actionAdminUsersDestroy(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&User{}, id).Error; err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/users")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("User was deleted.")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/users")
}

func actionAdminPagesIndex(c *gin.Context) {
	var pages []Page
	db.Find(&pages)
	c.HTML(http.StatusOK, "admin/pages/index.html", addFlashesAndUser(c, &gin.H{"pages": pages}))
}

func actionAdminPagesShow(c *gin.Context) {
	id := c.Param("id")
	var page Page

	if err := db.First(&page, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "admin/error.html", addFlashesAndUser(c, &gin.H{"errors": []string{"Page not found"}}))
		return
	}

	pageJSON, err := json.MarshalIndent(page, "", "  ")
	if err != nil {
		c.HTML(http.StatusNotFound, "admin/error.html", addFlashesAndUser(c, &gin.H{"errors": []string{err.Error()}}))
		return
	}

	c.HTML(http.StatusOK, "admin/pages/show.html", addFlashesAndUser(c, &gin.H{"page": page, "pageJSON": string(pageJSON)}))
}

func actionAdminPagesNew(c *gin.Context) {
	var page Page
	c.HTML(http.StatusOK, "admin/pages/new.html", addFlashesAndUser(c, &gin.H{"page": page}))
}

func actionAdminPagesCreate(c *gin.Context) {
	var page Page
	page.Slug = c.PostForm("slug")
	page.Content = c.PostForm("content")
	page.CreatedAt = time.Now()
	page.UpdatedAt = time.Now()

	page_input := &PageInput{
		Slug:    page.Slug,
		Content: page.Content,
	}

	// Validate user input
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(page_input); err != nil {
		c.HTML(http.StatusBadRequest, "admin/pages/new.html", addFlashesAndUser(c, &gin.H{"errors": humanValidationErrors(err), "page": page}))
		return
	}

	if err := db.Create(&page).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin/pages/new.html", addFlashesAndUser(c, &gin.H{"errors": []string{err.Error()}, "page": page}))
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Page was added.")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/pages")
}

func actionAdminPagesEdit(c *gin.Context) {
	id := c.Param("id")
	var page Page

	if err := db.First(&page, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "admin/error.html", addFlashesAndUser(c, &gin.H{"errors": []string{"User not found"}}))
		return
	}

	c.HTML(http.StatusOK, "admin/pages/edit.html", addFlashesAndUser(c, &gin.H{"page": page}))
}

func actionAdminPagesUpdate(c *gin.Context) {
	id := c.Param("id")
	var page Page

	if err := db.First(&page, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "admin/error.html", addFlashesAndUser(c, &gin.H{"errors": []string{"User not found"}}))
		return
	}

	page.Slug = c.PostForm("slug")
	page.Content = c.PostForm("content")
	page.UpdatedAt = time.Now()

	page_input := &PageInput{
		Slug:    page.Slug,
		Content: page.Content,
	}

	// Validate user input
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(page_input); err != nil {
		c.HTML(http.StatusOK, "admin/pages/edit.html", addFlashesAndUser(c, &gin.H{"errors": humanValidationErrors(err), "user": page}))
		return
	}

	if err := db.Save(&page).Error; err != nil {
		c.HTML(http.StatusOK, "admin/pages/edit.html", addFlashesAndUser(c, &gin.H{"errors": []string{err.Error()}, "user": page}))
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Page was edited.")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/pages")
}

func actionAdminPagesDestroy(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Page{}, id).Error; err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/pages")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Page was deleted.")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/pages")
}

func actionPublicTools(c *gin.Context) {
	c.HTML(http.StatusOK, "public/tools.html", addFlashesAndUser(c, &gin.H{}))
}

func actionPublicToolsDBClear(c *gin.Context) {
	session := sessions.Default(c)

	// Clear all users and pages from the database
	if err := db.Exec("delete from \"user\"").Error; err != nil {
		session.AddFlash("Error clearing users: " + err.Error())
	} else {
		if err := db.Exec("delete from \"page\"").Error; err != nil {
			session.AddFlash("Error clearing pages: " + err.Error())
		} else {
			session.AddFlash("Database cleared successfully.")
		}
	}

	session.Save()

	c.Redirect(http.StatusSeeOther, "/tools")
}

func actionPublicToolsSeed(c *gin.Context) {
	session := sessions.Default(c)

	result := db.Create(&User{Login: "admin", Password: "admin"})
	if result.Error != nil {
		session.AddFlash("Error seeding database: " + result.Error.Error())
	} else {
		result = db.Create(&Page{Slug: "about", Content: "This is the about page."})
		if result.Error != nil {
			session.AddFlash("Error seeding database: " + result.Error.Error())
		} else {
			session.AddFlash("Database seeded successfully.")
		}
	}

	session.Save()

	c.Redirect(http.StatusSeeOther, "/tools")
}

func actionPublicToolsSQL(c *gin.Context) {
	// Get SQL query from URL parameter
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'q' parameter"})
		return
	}

	sqlDB, err := db.DB()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting DB: " + err.Error()})
		return
	}

	rows, err := sqlDB.Query(q)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing SQL query: " + err.Error()})
		return
	}

	cols, err := rows.Columns()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving columns: " + err.Error()})
		return
	}

	out := make([]map[string]interface{}, 0)

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		out = append(out, m)
	}

	c.JSON(http.StatusOK, gin.H{"q": q, "out": out})
}
