package server

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/ichaf1997/alertmanager2/media"
	"github.com/ichaf1997/alertmanager2/utils"
	"github.com/sirupsen/logrus"
)

type Server struct {
	engine  *gin.Engine
	handler *media.Handler
}

func (server *Server) RegisterRoute() {

	server.engine.Use(gin.Recovery(), utils.StructuredLoggerHandlerFunc())

	server.engine.GET(
		"/_status/healthz",
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		},
	)

	v1 := server.engine.Group("/v1/medium")
	server.addFeiShuRoutes_v1(v1)

}

func (s *Server) Start(address string) error {
	s.RegisterRoute()
	logrus.Infof("Listening and serving HTTP on %s", address)
	return s.engine.Run(address)
}

func NewServer(address string, tmpl *template.Template) *Server {
	engine := gin.New()
	handler := media.NewHandler(tmpl)
	return &Server{engine: engine, handler: handler}
}
