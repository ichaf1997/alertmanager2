package cmd

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ichaf1997/alertmanager2/server"
	"github.com/ichaf1997/alertmanager2/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
			s := server.NewServer(listenAddress, utils.GetLocalTemplate(tmplDir))
			if err := s.Start(listenAddress); err != nil {
				logrus.Panicf("error starting alertmanager2 server: %v", err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&listenAddress, "web.listen-address", "a", "127.0.0.1:8080", "Addresses on which to expose metrics and web interface")
	serverCmd.Flags().StringVarP(&tmplDir, "tmpl.dirs", "t", filepath.Join(pwd, "templates", "*.tmpl"), "Go Templates directory to Load")
	serverCmd.MarkFlagRequired("tmpl.dir")
}
