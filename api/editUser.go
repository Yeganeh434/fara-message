package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
)

type Information struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func editUser(c *gin.Context) {
	var userInfo Information

	err := c.BindJSON(&userInfo)
	if err != nil {
		log.Printf("error binding JSON:%v", err)
		c.Status(400)
		return
	}

	if len(userInfo.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the username is incorrect",
		})
		return
	}

	userUnderReview, err := db.Mysql.ReadUserByUsername(userInfo.Username)
	if err != nil {
		log.Printf("error getting username from database:%v", err)
		c.Status(400)
		return
	}

	userUnderReview.FirstName = userInfo.Firstname
	userUnderReview.LastName = userInfo.Lastname
	err = db.Mysql.UpdateUser(userUnderReview.ID, db.ConvertUserToUserInfo(userUnderReview))
	if err != nil {
		log.Printf("error updating user:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "firstname and lastname edited successfully",
	})
}
