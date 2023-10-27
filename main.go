package main

import (
	"alertmanager2/channels"
	"alertmanager2/utils"
	"fmt"
	"net/http"

	"github.com/alecthomas/kingpin/v2"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
)

var (
	appName       = "alertmanager2"
	listenAddress = kingpin.Flag(
		"web.listen-address",
		"Address to listen on for UI, API, and telemetry.",
	).Default(":8080").String()
	enableWxwork = kingpin.Flag(
		"wxwork.enable-api",
		"Enable wxwork robot webhook mounted on /channel/wxwork path.",
	).Default("true").Bool()
	enableEs = kingpin.Flag(
		"es.enable-api",
		"Enable elastic-search webhook mounted on /channel/es path.",
	).Default("true").Bool()
	enableAli = kingpin.Flag(
		"ali.enable-api",
		"Enable alicloud webhook mounted on /channel/ali path.",
	).Default("true").Bool()
	runMode = kingpin.Flag(
		"mode",
		"Gin Framework mode. One of: [debug, release, test]",
	).Default(gin.ReleaseMode).String()
	loglevel = kingpin.Flag(
		"log.level",
		"Only log messages with the given severity or above. One of: [panic, fatal, error, warn, info, debug, trace]",
	).Default("info").String()
)

type routes struct {
	router *gin.Engine
}

func init() {

	kingpin.Version(version.Print(appName))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	gin.SetMode(*runMode)

}

func initRoute(logger *logrus.Logger) routes {
	r := routes{
		router: gin.New(),
	}

	r.router.GET(
		"/_status/healthz",
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		},
	)

	r.router.Use(gin.Recovery())
	r.router.Use(utils.LoggingMiddleware(logger))

	rg := r.router.Group("/channel")

	channels.SetGlobalLogger(logger)

	// Register channel routers
	if *enableAli {
		channels.NewChannelGroup(rg.Group("/ali")).Handle(channels.AliChannelRouters()...)
	}

	if *enableWxwork {
		channels.NewChannelGroup(rg.Group("/wxwork")).Handle(channels.WxWorkChannelRouters()...)
	}

	return r
}

func (r routes) Run(addr string, log *logrus.Logger) error {

	log.Info(fmt.Sprintf("Starting %s (%s - %s/%s)", appName, version.GoVersion, version.GoOS, version.GoArch))
	log.Info(version.Info())
	log.Info(fmt.Sprintf("Listening and serving HTTP on %s", *listenAddress))
	return r.router.Run(addr)
}

func main() {

	log := utils.DefaultLogger(*loglevel)
	route := initRoute(log)
	err := route.Run(*listenAddress, log)
	if err != nil {
		log.Error("error", err.Error())
	}

}
