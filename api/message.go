package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mhghw/fara-message/db"
)

func GetUserIDFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims format")
	}
	userID := claims[TokenUserID].(string)
	expirationTime := claims[TokenExpireTime].(time.Time)
	if expirationTime.Before(time.Now()) {
		return "", errors.New("token has expired")
	}
	return userID, nil
}

func GetUserID(authorizationHeader string) (string,error){
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "",errors.New("Unauthorized")
	}
	accessToken := parts[1]
	userID, err := GetUserIDFromToken(accessToken)
	if err != nil {
		return "",err
	}
	return userID,nil
}

type Message struct {
	ID       int    `json:"id"`
	SenderID int    `json:"senderID"`
	ChatID   int    `json:"chatID"`
	Content  string `json:"content"`
}

func DeleteMessageHandler(c *gin.Context) {
	var message Message
	if err := c.BindJSON(&message); err != nil {
		log.Printf("error binding JSON:%v", err)
	}
	err := db.Mysql.DeleteMessage(message.ID)
	if err != nil {
		log.Printf("error:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "message deleted successfully",
	})
}

func SendMessageHandler(c *gin.Context) {
	var message Message
	if err := c.BindJSON(&message); err != nil {
		log.Printf("error binding json:%v", err)
		c.Status(400)
		return
	}
	userID,err:=GetUserID(c.GetHeader("Authorization"))
	if err!=nil {
		log.Printf("error get user ID:%v",err)
		c.Status(400)
		return
	}
	ID, _ := strconv.Atoi(userID)
	message.SenderID = ID
	dbMessage:=db.Message{
		ID:message.ID,
		SenderID: message.SenderID,
		ChatID:message.ChatID,
		Content:message.Content,
	}
	err = db.Mysql.SendMessage(dbMessage)
	if err != nil {
		log.Printf("error:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "message sent successfully",
	})
}
