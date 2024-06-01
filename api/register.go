package api

import (
	// "encoding/json"
	"errors"
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
	if form.Username == "" || form.FirstName == "" || form.LastName == "" || form.Password == "" || form.DateOfBirth == "" || form.Email == "" {
		return errors.New("all fields must be filled")
	}

	if !(form.Gender == 0 || form.Gender == 1) {
		return errors.New("invalid gender")
	}

	if form.Password != form.ConfirmPassword {
		return errors.New("password does not match")
	}

	if !IsStrongPassword(form.Password) {
		return errors.New("your password must be at least 8 characters long and contain uppercase letter,lowercase letter,digit, and special character")
	}

	if len(form.Username) < 3 || len(form.Username) > 20 {
		return errors.New("username length should be between 3 and 20 characters")
	}
	spicificChar := "@#$%&*()+=!?,.<>/|~`\""
	for _, usernameValue := range form.Username {
		for _, charsValue := range spicificChar {
			if usernameValue == charsValue {
				return errors.New("username contains invalid characters")
			}
		}
	}
	isUsernameAvailable, err := db.Mysql.IsUsernameAvailable(form.Username)
	if err != nil {
		return err
	}
	if !isUsernameAvailable {
		return errors.New("this username is not available")
	}

	isEmailExist, err := db.Mysql.IsEmailExist(form.Email)
	if err != nil {
		return err
	}
	if isEmailExist {
		return errors.New("an account has already been created with this email")
	}
	if !isValidEmail(form.Email) {
		return errors.New("invalid email")
	}

	return nil
}

func convertRegisterFormToUser(form RegisterForm) (db.User, error) {
	layout := "2006-01-02"
	convertTime, err := time.Parse(layout, form.DateOfBirth)
	if err != nil {
		return db.User{}, err
	}

	id, err := generateID()
	if err != nil {
		return db.User{}, err
	}
	user := db.User{
		ID:          id,
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
