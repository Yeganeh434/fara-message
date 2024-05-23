package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
)

type Recipient struct {
	Email  string            `json:"Email"`
	Fields map[string]string `json:"Fields"`
}

type Body struct {
	ContentType string `json:"ContentType"`
	Content     string `json:"Content"`
	Charset     string `json:"Charset"`
}

type Content struct {
	Body    []Body `json:"Body"`
	Subject string `json:"Subject"`
	From    string `json:"From"`
}

type EmailRequest struct {
	Recipients []Recipient `json:"Recipients"`
	Content    Content     `json:"Content"`
}

func GetOTP(c *gin.Context) {
	otp := GenerateOTP()
	email := c.Param("email")
	isEmailExist, err := db.Mysql.IsEmailExist(email)
	if err != nil {
		log.Printf("error in checking the existence of email:%v", err)
		c.Status(400)
		return
	}
	if !isEmailExist {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "there is no user with this email",
		})
	}
	emailRequest := EmailRequest{
		Recipients: []Recipient{
			{
				Email: email,
				Fields: map[string]string{
					"otp": otp,
				},
			},
		},
		Content: Content{
			Body: []Body{
				{
					ContentType: "HTML",
					Content:     "<p>Your OTP code is: {otp}</p>",
					Charset:     "utf-8",
				},
			},
			Subject: "change password",
			From:    "toncas <toncas224@gmail.com>",
		},
	}

	requestBody, err := json.Marshal(emailRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling request body"})
		return
	}

	req, err := http.NewRequest("POST", "https://api.elasticemail.com/v4/emails", bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ElasticEmail-ApiKey", "A328AE6824853599EA548A299C42C6FAF06C06BDF0BD5569BEED4401F03E06E3D4B59D1A39792E1AF67733FA2116A428")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request"})
		return
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}

	//save otp in database
	id, _ := strconv.Atoi(generateID())
	IntOfOTP, _ := strconv.Atoi(otp)
	otpInfo := db.OTP{
		ID:    id,
		OTP:   IntOfOTP,
		Email: email,
	}
	err = db.Mysql.SaveOTPInDB(otpInfo)
	if err != nil {
		log.Printf("error in saving OTP in the database:%v", err)
		c.Status(400)
		return
	}

	c.JSON(resp.StatusCode, gin.H{
		"headers": resp.Header,
		"body":    string(responseBody),
	})
}

// change this func!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!1
func GenerateOTP() string {
	const charset = "0123456789"
	rand.NewSource(10)
	id := make([]byte, 5)
	for idx := range id {
		id[idx] = charset[rand.Intn(len(charset))]
	}

	return string(id)
}

type newPasswordRequest struct {
	OTP             int    `json:"otp"`
	Email           string `json:"email"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

func ChangePasswordHandler(c *gin.Context) {
	var newPasswordInfo newPasswordRequest
	err := c.BindJSON(&newPasswordInfo)
	if err != nil {
		log.Printf("error binding JSON:%v", err)
		c.Status(400)
		return
	}
	isOTPCorrect, err := db.Mysql.IsOTPCorrect(newPasswordInfo.OTP, newPasswordInfo.Email)
	if err != nil {
		log.Printf("error checking whether OTP is correct:%v", err)
		c.Status(400)
		return
	}
	if !isOTPCorrect {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the OTP you entered is not correct",
		})
	}

	if newPasswordInfo.NewPassword != newPasswordInfo.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "the password does not match the password confirmation",
		})
		return
	}
	if !IsStrongPassword(newPasswordInfo.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "your password must be at least 8 characters long and contain uppercase letter,lowercase letter,digit, and special character",
		})
		return
	}
	newPasswordInfo.NewPassword = hash(newPasswordInfo.NewPassword)
	err = db.Mysql.ChangePassword(newPasswordInfo.Email, newPasswordInfo.NewPassword)
	if err != nil {
		log.Printf("error changing password:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "password change successfully",
	})
}
