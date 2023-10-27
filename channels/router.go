package channels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type Router struct {
	HandlerMethod string
	HandlerPath   string
	HandlerFunc   gin.HandlerFunc
}

type ChannelGroup struct {
	group *gin.RouterGroup
}

func SetGlobalLogger(logger *logrus.Logger) {
	log = logger
}

func NewChannelGroup(cg *gin.RouterGroup) ChannelGroup {
	return ChannelGroup{
		cg,
	}
}

func (cg ChannelGroup) Handle(rs ...Router) {
	for _, ch := range rs {
		cg.group.Handle(
			ch.HandlerMethod,
			ch.HandlerPath,
			ch.HandlerFunc,
		)
	}
}

func AliChannelRouters() []Router {
	var rs = []Router{
		{
			HandlerMethod: http.MethodPost,
			HandlerPath:   "/sendsms",
			HandlerFunc:   SendSms,
		},
		{
			HandlerMethod: http.MethodPost,
			HandlerPath:   "/sendvms",
			HandlerFunc:   SendVms,
		},
		{
			HandlerMethod: http.MethodPost,
			HandlerPath:   "/sendsimplevms",
			HandlerFunc:   SendFastVms,
		},
		{
			HandlerMethod: http.MethodPost,
			HandlerPath:   "/sendsimplesms",
			HandlerFunc:   SendFastSms,
		},
	}
	return rs
}

func WxWorkChannelRouters() []Router {
	var rs = []Router{
		{
			HandlerMethod: http.MethodPost,
			HandlerPath:   "/robot",
			HandlerFunc:   SendWxWorkRobot,
		},
	}
	return rs
}
