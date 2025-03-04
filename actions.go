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

type UserInput struct {
	Login    string `validate:"required,min=3"`
	Password string `validate:"required,min=3"`
}

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
		return nil
	}

	return h
}

func actionRoot(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func actionLoginForm(c *gin.Context) {
	_, exists := c.Get("currentUser")
	if exists {
		c.Redirect(http.StatusSeeOther, "/admin/users")
	}

	c.HTML(http.StatusOK, "login.html", nil)
}

func actionLoginSubmit(c *gin.Context) {
	_, exists := c.Get("currentUser")
	if exists {
		c.Redirect(http.StatusSeeOther, "/admin/users")
	}

	username := c.PostForm("login")
	password := c.PostForm("password")

	var user User
	// TODO: Encrypt password
	if err := db.Where("login = ? and password = ?", username, password).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"errors": []string{"Invalid username or password"}})
		return
	}

	session := sessions.Default(c)
	session.Set("currentUser", user.ID)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/users")
}

func actionLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("currentUser")
	session.AddFlash("Logged out")
	session.Save()
	c.Redirect(http.StatusSeeOther, "/admin/users")
}

func actionUsersIndex(c *gin.Context) {
	var users []User
	db.Find(&users)
	c.HTML(http.StatusOK, "users_index.html", addFlashesAndUser(c, &gin.H{"users": users}))
}

func actionUsersShow(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := db.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", addFlashesAndUser(c, &gin.H{"errors": []string{"User not found"}}))
		return
	}

	userJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", addFlashesAndUser(c, &gin.H{"errors": []string{err.Error()}}))
		return
	}

	c.HTML(http.StatusOK, "users_show.html", addFlashesAndUser(c, &gin.H{"user": user, "userJSON": string(userJSON)}))
}

func actionAdminIndex(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/admin/users")
}

func actionUsersNew(c *gin.Context) {
	var user User
	c.HTML(http.StatusOK, "users_new.html", addFlashesAndUser(c, &gin.H{"user": user}))
}

func actionUsersCreate(c *gin.Context) {
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
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(user_input); err != nil {
		c.HTML(http.StatusBadRequest, "users_new.html", addFlashesAndUser(c, &gin.H{"errors": humanValidationErrors(err), "user": user}))
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "users_new.html", addFlashesAndUser(c, &gin.H{"errors": []string{err.Error()}, "user": user}))
		return
	}

	session := sessions.Default(c)
	session.AddFlash("User was added.")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/users")
}

func actionUsersEdit(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := db.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", addFlashesAndUser(c, &gin.H{"errors": []string{"User not found"}}))
		return
	}

	c.HTML(http.StatusOK, "users_edit.html", addFlashesAndUser(c, &gin.H{"user": user}))
}

func actionUsersUpdate(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := db.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", addFlashesAndUser(c, &gin.H{"errors": []string{"User not found"}}))
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
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(user_input); err != nil {
		c.HTML(http.StatusOK, "users_edit.html", addFlashesAndUser(c, &gin.H{"errors": humanValidationErrors(err), "user": user}))
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.HTML(http.StatusOK, "users_edit.html", addFlashesAndUser(c, &gin.H{"errors": []string{err.Error()}, "user": user}))
		return
	}

	session := sessions.Default(c)
	session.AddFlash("User was edited.")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/users")
}

func actionUsersDestroy(c *gin.Context) {
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
