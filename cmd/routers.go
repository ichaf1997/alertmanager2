package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/ichaf1997/alertmanager2/media"
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

func (s *Server) addBytesRoutes_v1(rg *gin.RouterGroup) {

	r := rg.Group("/feishu")

	r.POST("/test", media.HelloWorld)
	r.POST("/robot/:key", s.chs.FeiShuCustomRobotHandler())

}
