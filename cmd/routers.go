package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/ichaf1997/alertmanager2/channels"
)

// func addAliCloudRoutes_v1(rg *gin.RouterGroup) {

// 	r := rg.Group("/ali")

// 	r.POST("/sendsms", channels.SendSms)
// 	r.POST("/sendvms", channels.SendVms)
// 	r.POST("/sendsimplesms", channels.SendFastSms)
// 	r.POST("/sendsimplevms", channels.SendFastVms)

// }

// func addWxworkRoutes_v1(rg *gin.RouterGroup) {

// 	r := rg.Group("/wxwork")

// 	r.POST("/robot", channels.SendWxWorkRobot)

// }

func addBytesRoutes_v1(rg *gin.RouterGroup) {

	r := rg.Group("/feishu")

	r.POST("/robot", channels.HelloWorld)

}
