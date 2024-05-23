package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func RunWebServer(port int) error {
	addr := fmt.Sprintf(":%d", port)
	router := gin.New()
	router.POST("/user/register", RegisterHandler)
	router.GET("/user/get_otp/:email", GetOTPHandler)
	router.POST("/user/change_password", ChangePasswordHandler)
	router.POST("/login", authenticateUser)
	router.Use(AuthMiddlewareHandler)
	router.POST("/user/read", ReadUserHandler)
	router.POST("/user/update", UpdateUserHandler)
	router.POST("/user/delete", DeleteUserHandler)
	router.POST("/user/contact/:contactID",AddContactHandler)
	router.DELETE("/user/contact/:contactID",DeleteContactHandler)
	router.GET("/user/contact",GetContactHandler)
	router.POST("/send/message", SendMessageHandler)
	router.DELETE("/delete/message/:id", DeleteMessageHandler)
	router.POST("/new_direct_chat", NewDirectChatHandler)
	router.POST("/new_group_chat", NewGroupChatHandler)
	router.POST("/chat/group/add",AddMemberToGroupHandler)
	router.GET("/chat/messages/:id", GetChatMessagesHandler)
	router.GET("/user/chat/members/:id",GetChatMembersHandler)
	router.GET("/user/chat/list",GetChatsListHandler)
	err := router.Run(addr)
	return err
}
