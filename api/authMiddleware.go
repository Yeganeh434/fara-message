package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
)

func AuthMiddlewareHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusForbidden, gin.H{

			"error": "Authorization is required",
		})

		c.AbortWithStatus(http.StatusForbidden)

		return
	}
	userID, err := ValidateToken(tokenString)
	if err != nil {
		log.Printf("failed to validate token:%v",err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	_, err = db.Mysql.ReadUser(userID)
	if err != nil {
		log.Print("user ID is not in the DataBase: %w", err)
		c.AbortWithStatus(http.StatusForbidden)

		return
	}

	c.Next()
}
