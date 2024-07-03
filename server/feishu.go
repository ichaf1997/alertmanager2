package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ichaf1997/alertmanager2/media"
)

func (s *Server) addFeiShuRoutes_v1(rg *gin.RouterGroup) {

	feishu := rg.Group("/feishu")

	feishu.POST("/test", media.HelloWorld)
	feishu.POST("/robot/:key", s.handler.FeiShuCustomRobotHandler())

}
