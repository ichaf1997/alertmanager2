package channels

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type FeiShuCustomRobot struct {
	webhook string
	// msgData     map[string]interface{}
	// msgTmplName string
	// secrets     []string
}

func HelloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func (ch *Channel) FeiShuCustomRobotFunc() func(c *gin.Context) {
	return func(c *gin.Context) {
		var f FeiShuCustomRobot
		c.Bind(&f)
		fmt.Println(f.webhook)
	}
}
