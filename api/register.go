package api

import (
	// "encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
)

type RegisterForm struct {
	Username        string `json:"username"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Gender          int    `json:"gender"`
	DateOfBirth     string `json:"dateOfBirth"`
	Email           string `json:"email"`
}

type tokenJSON struct {
	Token string `json:"token"`
}

func RegisterHandler(c *gin.Context) {
	var requestBody RegisterForm
	err := c.BindJSON(&requestBody)
	if err != nil {
		log.Printf("failed to bind json:%v", err)
		c.Status(400)
		return
	}
	err = validateUser(requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"failed to validate user": err.Error(),
		})
		return
	}

	requestBody.Password = hash(requestBody.Password)
	user, err := convertRegisterFormToUser(requestBody)
	if err != nil {
		log.Printf("failed to convert register form to user:%v", err)
		c.Status(400)
		return
	}

	token, err := CreateJWTToken(user.ID)

	if err != nil {
		log.Print("failed to create token")
		c.Status(400)
		return
	}
	userToken := tokenJSON{
		Token: token,
	}

	db.Mysql.CreateUser(user)
	c.JSON(http.StatusOK, userToken)
}

// other validation fields will be added...
func validateUser(form RegisterForm) error {
	if form.Password != form.ConfirmPassword {
		return errors.New("password does not match")
	}

	if !IsStrongPassword(form.Password) {
		return errors.New("your password must be at least 8 characters long and contain uppercase letter,lowercase letter,digit, and special character")
	}

	isUsernameAvailable, err := db.Mysql.IsUsernameAvailable(form.Username)
	if err != nil {
		return err
	}
	if !isUsernameAvailable {
		return errors.New("this username is not available")
	}

	return nil
}

func convertRegisterFormToUser(form RegisterForm) (db.User, error) {
	layout := "2006-01-02"
	convertTime, err := time.Parse(layout, form.DateOfBirth)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to parse date %w", err)
	}

	generatedID := generateID()
	user := db.User{
		ID:          generatedID,
		Username:    form.Username,
		FirstName:   form.FirstName,
		LastName:    form.LastName,
		Password:    form.Password,
		Gender:      form.Gender,
		DateOfBirth: convertTime,
		CreatedTime: time.Now(),
		Email:       form.Email,
	}

	return user, nil
}
