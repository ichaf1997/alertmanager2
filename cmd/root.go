package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "alertmanager2",
		Short: "Alertmanager Webhook with integration 'AliSms' 'AliVms' 'WxworkRobot' 'FeiShuRobot' 'ElasticSearchAPI' etc.",
	}
	verbose  bool
	loglevel string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show verbose log")
	rootCmd.PersistentFlags().StringVar(&loglevel, "loglevel", "info", "Only log messages with the given severity or above. One of: [error, warn, info, debug]")
}
