package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
)

type loginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func authenticateUser(c *gin.Context) {
	var loginData loginBody
	err := c.BindJSON(&loginData)
	if err != nil {
		log.Printf("error binding JSON:%v", err)
		c.Status(400)
		return
	}

	errIncorrectUserOrPass := HTTPError{Message: "the username or password is incorrect"}
	errIncorrectUserOrPassJSON, errInMarshalling := json.Marshal(errIncorrectUserOrPass)
	if errInMarshalling != nil {
		log.Printf("error:%v", errInMarshalling)
		c.Status(400)
		return
	}

	if len(loginData.Username) < 3 || len(loginData.Password) < 8 {
		c.JSON(http.StatusBadRequest, errIncorrectUserOrPassJSON)
		return
	}

	//checking entered data with data that is already stored
	userUnderReveiw, err := db.Mysql.ReadUserByUsername(loginData.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, errIncorrectUserOrPassJSON)
	}
	if hash(loginData.Password)!=userUnderReveiw.Password {
		log.Printf("the password is incorrect")
		c.Status(400)
		return
	}

	token, err := CreateJWTToken(userUnderReveiw.ID)
	if err != nil {
		log.Print("failed to create token")
		return
	}
	userToken := tokenJSON{
		Token: token,
	}
	c.JSON(http.StatusOK,userToken)
}
