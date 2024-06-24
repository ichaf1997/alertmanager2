package cmd

import (
	"net/http"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/ichaf1997/alertmanager2/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Server struct {
	engine        *gin.Engine
	localTemplate *template.Template
}

var (
	listenAddress string
	tmplDir       string
	serverCmd     = &cobra.Command{
		Use:   "server",
		Short: "Start the alertmanager2 server",
		Run: func(cmd *cobra.Command, args []string) {
			utils.InitLogger(loglevel, verbose)
			if logrus.GetLevel() != logrus.DebugLevel {
				gin.SetMode(gin.ReleaseMode)
			}
			server := Server{
				engine:        gin.New(),
				localTemplate: utils.GetLocalTemplate(tmplDir),
			}
			if err := server.Start(listenAddress); err != nil {
				logrus.Panicf("error starting alertmanager2 server: %v", err)
				os.Exit(1)
			}
		},
	}
)

func (server Server) RegisterRoute() {

	server.engine.Use(gin.Recovery(), utils.StructuredLoggerHandlerFunc())

	server.engine.GET(
		"/_status/healthz",
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		},
	)

	v1 := server.engine.Group("/v1/channel")
	// addAliCloudRoutes_v1(v1)
	// addWxworkRoutes_v1(v1)
	addBytesRoutes_v1(v1)

}

func (server Server) Start(address string) error {
	server.RegisterRoute()
	logrus.Infof("Listening and serving HTTP on %s", address)
	return server.engine.Run(address)
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&listenAddress, "web.listen-address", "a", "127.0.0.1:8080", "Addresses on which to expose metrics and web interface")
	serverCmd.Flags().StringVarP(&tmplDir, "tmpl.dirs", "t", "", "Go Templates directory to Load")
	serverCmd.MarkFlagRequired("tmpl.dir")
}
