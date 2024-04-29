package api

import (
	"encoding/json"
	"log"
	"net/http"

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
	if c.Request.ContentLength != 0 {
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
		userID,err:=GetUserID(c.GetHeader("Authorization"))
		if err!=nil {
			log.Printf("error get user ID:%v",err)
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
	userID,err:=GetUserID(c.GetHeader("Authorization"))
	if err!=nil {
		log.Printf("error get user ID:%v",err)
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
	dbUserInfo:=ConvertUserInfo(newInfo)
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
	userID,err:=GetUserID(c.GetHeader("Authorization"))
	if err!=nil {
		log.Printf("error get user ID:%v",err)
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
