package api

import (
	"fmt"
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
	messageIDString := c.Param("id")
	messageID, err := strconv.Atoi(messageIDString)
	if err != nil {
		log.Printf("error converting message Id to int:%v", err)
		c.Status(400)
		return
	}
	userIDString, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("error get user ID:%v", err)
		c.Status(400)
		return
	}
	userID, _ := strconv.Atoi(userIDString)
	var message db.Message
	message, err = db.Mysql.FindMessageInfo(messageID)
	if err != nil {
		log.Printf("error in finding message information:%v", err)
		c.Status(400)
		return
	}

	//anyone can just delete their own message
	// if message.SenderID!=userID {
	// 	c.JSON(http.StatusBadRequest,gin.H{
	// 		"error":"you are not allowed to delete this message",
	// 	})
	// }

	//any member of the chat can delete any message from the chat
	isChatContact, err := db.Mysql.IsAChatContact(userID, message.ChatID)
	if err != nil {
		log.Printf("error in checking the existence of a contact in the chat:%v", err)
		c.Status(400)
		return
	}
	if !isChatContact {
		fmt.Println(message.SenderID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you are not allowed to delete this message",
		})
		return
	}
	err = db.Mysql.DeleteMessage(message)
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
	senderID, _ := strconv.Atoi(userID)
	messageID, _ := strconv.Atoi(generateID())
	dbMessage := db.Message{
		ID:       messageID,
		SenderID: senderID,
		ChatID:   message.ChatID,
		Content:  message.Content,
	}
	isChatContact, err := db.Mysql.IsAChatContact(dbMessage.SenderID, dbMessage.ChatID)
	if err != nil {
		log.Printf("error in checking the existence of a contact in the chat:%v", err)
		c.Status(400)
		return
	}
	if !isChatContact {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you are not allowed to send messages in this chat",
		})
		return
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
