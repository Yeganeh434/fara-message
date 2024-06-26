package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
)

type UsernameType struct {
	Username string `json:"username"`
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
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		var username UsernameType
		err := c.BindJSON(&username)
		if err != nil {
			log.Printf("error binding json:%v", err)
			c.Status(400)
			return
		}
		user, err := db.Mysql.ReadAnotherUser(username.Username)
		if err != nil {
			log.Printf("error reading user:%v", err)
			c.Status(400)
			return
		} else {
			var userInfo AnotherUserInfo
			userInfo.Username = user.Username
			userInfo.FirstName = user.FirstName
			userInfo.LastName = user.LastName
			convertUserToJSON, err := json.Marshal(userInfo)
			if err != nil {
				log.Printf("error marshaling:%v", err)
				c.Status(400)
				return
			}
			c.JSON(http.StatusOK, convertUserToJSON)
		}
	} else {
		parts := strings.Split(authorizationHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
			return
		}
		accessToken := parts[1]
		userID, err := getUserIDFromToken(accessToken)
		if err != nil {
			log.Printf("error:%v", err)
			c.Status(400)
			return
		}
		user, err := db.Mysql.ReadUser(userID)
		if err != nil {
			log.Printf("error reading user:%v", err)
			c.Status(400)
			return
		} else {
			var userInfo UserInfo

			userInfo.Username = user.Username
			userInfo.FirstName = user.FirstName
			userInfo.LastName = user.LastName
			userInfo.Gender = int(user.Gender)
			userInfo.DateOfBirth = user.DateOfBirth
			userInfo.CreatedTime = user.CreatedTime
			convertUserToJSON, err := json.Marshal(userInfo)
			if err != nil {
				log.Printf("error marshaling:%v", err)
				c.Status(400)
				return
			}
			c.JSON(http.StatusOK, convertUserToJSON)
		}
	}
}

func UpdateUserHandler(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	accessToken := parts[1]
	userID, err := getUserIDFromToken(accessToken)
	if err != nil {
		log.Printf("error:%v", err)
		c.Status(400)
		return
	}
	var newInfo UserInfo
	err = c.BindJSON(&newInfo)
	if err != nil {
		log.Printf("error binding JSON:%v", err)
		c.Status(400)
		return
	}
	gender := db.Male
	if newInfo.Gender != 0 {
		gender = db.Female
	}

	dbUserInfo := db.UserInfo{
		Username:    newInfo.Username,
		FirstName:   newInfo.FirstName,
		LastName:    newInfo.LastName,
		Gender:      gender,
		DateOfBirth: newInfo.DateOfBirth,
		CreatedTime: newInfo.CreatedTime,
	}
	err = db.Mysql.UpdateUser(userID, dbUserInfo)
	if err != nil {
		log.Printf("error updating user:%v", err)
		c.Status(400)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "user updated successfully",
		})
	}
}

func DeleteUserHandler(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	accessToken := parts[1]
	userID, err := getUserIDFromToken(accessToken)
	if err != nil {
		log.Printf("error:%v", err)
		c.Status(400)
		return
	}
	err = db.Mysql.DeleteUser(userID)
	if err != nil {
		log.Printf("error:%v", err)
		c.Status(400)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "user deleted successfully",
		})
	}
}
