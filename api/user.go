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
}

// func CreateUserHandler(c *gin.Context) {
// 	var newUser db.User
// 	err := c.BindJSON(&newUser)
// 	if err != nil {
// 		log.Printf("error binding JSON:%v", err)
// 		c.Status(400)
// 		return
// 	}
// 	err = db.CreateUser(newUser)
// 	if err != nil {
// 		log.Printf("error inserting user:%v", err)
// 		c.Status(400)
// 		return
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "user inserted successfully",
// 		})
// 	}

// }

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

func changePasswordHandler(c *gin.Context) {
	userID, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("error get user ID:%v", err)
		c.Status(400)
		return
	}
	var newPassword PasswordType
	err = c.BindJSON(&newPassword)
	if err != nil {
		log.Printf("error binding JSON:%v", err)
		c.Status(400)
		return
	}
	if len(newPassword.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "password is too short!",
		})
		return
	}
	if newPassword.Password != newPassword.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "the password does not match the password confirmation",
		})
		return
	}
	err = db.Mysql.ChangePassword(userID, newPassword.Password)
	if err != nil {
		log.Printf("error changing password:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "password change successfully",
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
	err = db.Mysql.AddContact(intOfUserID, intOfContactID)
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
	var contacts ContactsResponse
	dbContacts, err := db.Mysql.GetContact(intOfUserID)
	if err != nil {
		log.Printf("error get contacts:%v", err)
		c.Status(400)
		return
	}
	for i, value := range dbContacts {
		contacts.Contacts[i].Username = value.Username
		contacts.Contacts[i].FirstName = value.FirstName
		contacts.Contacts[i].LastName = value.LastName
	}
	c.JSON(http.StatusOK, contacts)
}
