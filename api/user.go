package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
)

type UsernameType struct {
	Username string `json:"username"`
}

type PasswordType struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type UpdateUser struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Gender      int    `json:"gender"`
	DateOfBirth string `json:"dateOfBirth"`
	Email       string `json:"email"`
}

func ReadUserHandler(c *gin.Context) {
	if c.Request.ContentLength != 0 {
		var username UsernameType
		err := c.BindJSON(&username)
		if err != nil {
			log.Printf("error binding json:%v", err)
			c.Status(400)
			return
		}
		user, err := db.Mysql.ReadUserByUsername(username.Username)
		if err != nil {
			log.Printf("error reading user:%v", err)
			c.Status(400)
			return
		} else {
			var userInfo AnotherUserInfo
			userInfo.Username = user.Username
			userInfo.FirstName = user.FirstName
			userInfo.LastName = user.LastName
			c.JSON(http.StatusOK, userInfo)
		}
	} else {
		userID, err := GetUserID(c.GetHeader("Authorization"))
		if err != nil {
			log.Printf("error get user ID:%v", err)
			c.Status(400)
			return
		}
		user, err := db.Mysql.ReadUser(userID)
		if err != nil {
			log.Printf("error reading user:%v", err)
			c.Status(400)
			return
		} else {
			userInfo := ConvertdbUser(user)
			c.JSON(http.StatusOK, userInfo)
		}
	}
}

func UpdateUserHandler(c *gin.Context) {
	userID, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("error get user ID:%v", err)
		c.Status(400)
		return
	}
	var newInfo UpdateUser
	err = c.BindJSON(&newInfo)
	if err != nil {
		log.Printf("error binding JSON:%v", err)
		c.Status(400)
		return
	}
	dbUserInfo := ConvertUpdateUser(newInfo)
	err = db.Mysql.UpdateUser(userID, dbUserInfo)
	if err != nil {
		log.Printf("error updating user:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
	})
}

func DeleteUserHandler(c *gin.Context) {
	userID, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("error get user ID:%v", err)
		c.Status(400)
		return
	}
	err = db.Mysql.DeleteUser(userID)
	if err != nil {
		log.Printf("error deleting user:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})
}


func AddContactHandler(c *gin.Context) {
	userID, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("error get user ID:%v", err)
		c.Status(400)
		return
	}
	contactID := c.Param("contactID")
	if contactID == "" {
		log.Print("contact ID is empty")
		c.Status(400)
		return
	}
	intOfUserID, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("error converting to int:%v", err)
		c.Status(400)
		return
	}
	intOfContactID, err := strconv.Atoi(contactID)
	if err != nil {
		log.Printf("error converting to int:%v", err)
		c.Status(400)
		return
	}
	ID, err := strconv.Atoi(generateID())
	if err != nil {
		log.Printf("error converting to int:%v", err)
		c.Status(400)
		return
	}
	isContactExist, err := db.Mysql.IsContactExist(intOfUserID, intOfContactID)
	if err != nil {
		log.Printf("error:%v", err)
		c.Status(400)
		return
	}
	if isContactExist {
		c.JSON(http.StatusOK, gin.H{
			"message": "you have already added this contact",
		})
		return
	}
	err = db.Mysql.AddContact(ID, intOfUserID, intOfContactID)
	if err != nil {
		log.Printf("error add contact to database:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "contact successfully added",
	})
}

func DeleteContactHandler(c *gin.Context) {
	userID, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("error get user ID:%v", err)
		c.Status(400)
		return
	}
	contactID := c.Param("contactID")
	if contactID == "" {
		log.Print("contact ID is empty")
		c.Status(400)
		return
	}
	intOfUserID, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("error converting to int:%v", err)
		c.Status(400)
		return
	}
	intOfContactID, err := strconv.Atoi(contactID)
	if err != nil {
		log.Printf("error converting to int:%v", err)
		c.Status(400)
		return
	}
	err = db.Mysql.DeleteContact(intOfUserID, intOfContactID)
	if err != nil {
		log.Printf("error deleting contact:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "contact deleted successfully",
	})
}

func GetContactHandler(c *gin.Context) {
	userID, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("error get user ID:%v", err)
		c.Status(400)
		return
	}
	intOfUserID, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("error converting to int:%v", err)
		c.Status(400)
		return
	}
	dbContacts, err := db.Mysql.GetContact(intOfUserID)
	if err != nil {
		log.Printf("error get contacts:%v", err)
		c.Status(400)
		return
	}
	var contacts []AnotherUserInfo
	for _, value := range dbContacts {
		contacts = append(contacts, AnotherUserInfo{
			Username:  value.Username,
			FirstName: value.FirstName,
			LastName:  value.LastName,
		})
	}
	var contactsResponse ContactsResponse
	for _, value := range contacts {
		contactsResponse.Contacts = append(contactsResponse.Contacts, value)
	}
	c.JSON(http.StatusOK, contactsResponse)
}
