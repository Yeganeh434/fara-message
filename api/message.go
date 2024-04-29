package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
)

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
	userID, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("error get user ID:%v", err)
		c.Status(400)
		return
	}
	ID, _ := strconv.Atoi(userID)
	message.SenderID = ID
	dbMessage := db.Message{
		ID:       message.ID,
		SenderID: message.SenderID,
		ChatID:   message.ChatID,
		Content:  message.Content,
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
