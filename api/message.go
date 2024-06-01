package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
)

type Message struct {
	ID       uint64      `json:"id"`
	SenderID uint64       `json:"senderID"`
	ChatID   uint64       `json:"chatID"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
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
	userID, _ := strconv.ParseUint(userIDString,10,64)
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
	senderID, _ := strconv.ParseUint(userID,10,64)
	messageID, err := generateID()
	if err != nil {
		log.Printf("error in generating ID:%v", err)
		c.Status(400)
		return
	}
	dbMessage := db.Message{
		ID:       messageID,
		SenderID: senderID,
		ChatID:   message.ChatID,
		Content:  message.Content,
		Time:     time.Now(),
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
