package media

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ichaf1997/alertmanager2/utils"
	"github.com/sirupsen/logrus"
)

type FeiShuCustomRobot struct {
	MsgTmplName string `form:"msg_tmpl_name" binding:"required"`
	MsgData     any
	Key         string
}

type FeiShuCustomRobotCard struct {
	MsgType string         `json:"msg_type"`
	Card    map[string]any `json:"card"`
}

// func (feishu FeiShuCustomRobot) SendMsg() error {
// 	// TODO:
// }

func (feishu FeiShuCustomRobot) Info() any {
	
}

func HelloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func (ch *Channel) FeiShuCustomRobotHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			feishu FeiShuCustomRobot
			data   utils.Notification
		)

		if err1, err2 := c.BindQuery(&feishu), c.BindJSON(&data); err1 == nil && err2 == nil {
			feishu.Key = fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", c.Param("key"))
			// c.JSON(200, feishu.Key)
			// c.JSON(200, feishu)
			// c.JSON(200, data)
		} else if err1 != nil {
			logrus.Warnf("Failed to parse query: %v", err1)
		} else if err2 != nil {
			logrus.Warnf("Failed to parse post data: %v", err2)
		}

	}
}
