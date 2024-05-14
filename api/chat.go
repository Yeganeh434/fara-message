package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
)

type GroupChatRequest struct {
	ChatName string   `json:"chatName"`
	Users    []string `json:"users"`
}
type DirectChatRequest struct {
	User string `json:"user"`
}

func NewDirectChatHandler(c *gin.Context) {
	user1, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("error get user ID:%v", err)
		c.Status(400)
		return
	}

	var requestBody DirectChatRequest
	err = c.BindJSON(&requestBody)
	if err != nil {
		log.Printf("failed to bind json,%v ", err)
		c.Status(400)
		return
	}
	user2, err := db.Mysql.ReadUserByUsername(requestBody.User)
	if err != nil {
		log.Printf("error getting user:%v ", err)
		c.Status(400)
		return
	}

	var users []string
	intOfUser1, _ := strconv.Atoi(user1)
	intOfUser2, _ := strconv.Atoi(user2.ID)
	if intOfUser1 < intOfUser2 {
		users = append(users, user1, user2.ID)
	} else {
		users = append(users, user2.ID, user1)
	}
	isChatExist, err := db.Mysql.IsExist(users)
	if err != nil {
		log.Printf("error:%v ", err)
		c.Status(400)
		return
	}
	if isChatExist {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "this chat has already been created",
		})
		return
	}
	if err := db.Mysql.NewChat(generateID(), "", 0, users); err != nil {
		log.Print("failed to create chat, ", err)
		c.Status(400)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "direct chat created successfully",
	})
}

func NewGroupChatHandler(c *gin.Context) {
	user1, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("error get user ID:%v", err)
		c.Status(400)
		return
	}
	var users []string
	users = append(users, user1)

	var requestBody GroupChatRequest
	err = c.BindJSON(&requestBody)
	if err != nil {
		log.Print("failed to bind json, ", err)
		return
	}
	for _, value := range requestBody.Users {
		user, err := db.Mysql.ReadUserByUsername(value)
		if err != nil {
			log.Printf("error getting user:%v ", err)
			return
		}
		users = append(users, user.ID)
	}

	if err := db.Mysql.NewChat(generateID(), requestBody.ChatName, 1, users); err != nil {
		log.Printf("failed to create chat: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "group chat created successfully",
	})
}

func GetChatMessagesHandler(c *gin.Context) {
	chatIDString := c.Param("id")
	chatID, _ := strconv.Atoi(chatIDString)
	messages, err := db.Mysql.GetChatMessages(chatID)
	if err != nil {
		log.Print(err)
		return
	}
	c.JSON(200, messages)
}

// response!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func GetChatMembersHandler(c *gin.Context) {
	chatIDString := c.Param("id")
	chatID, _ := strconv.Atoi(chatIDString)
	chatMembers, err := db.Mysql.GetChatMembers(chatID)
	if err != nil {
		log.Print("failed to get users chat members")
		return
	}
	c.JSON(200, chatMembers)
}

func GetChatsListHandler(c *gin.Context) {
	userID, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("error get user ID:%v", err)
		c.Status(400)
		return
	}
	chatsName, err := db.Mysql.GetChatsList(userID)
	if err != nil {
		log.Printf("error get chats name:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, chatsName)
}
