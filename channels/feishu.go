package channels

import (
	"github.com/gin-gonic/gin"
)

type Channel interface {
	sendNotification(c *gin.Context)
}

func HelloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}
