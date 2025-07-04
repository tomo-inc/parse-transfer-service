package service

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tomo.inc/parse-transfer-service/app/handler"
	"github.com/tomo.inc/parse-transfer-service/app/service"
	"github.com/tomo.inc/parse-transfer-service/app/service/chain"
	"github.com/tomo.inc/parse-transfer-service/cmd/service/config"
	"github.com/tomo.inc/parse-transfer-service/pkg/bot"
	"github.com/tomo.inc/parse-transfer-service/share"
	"time"
)

var cfgFile string

var serviceCmd = &cobra.Command{
	Use: "service",
	RunE: func(cmd *cobra.Command, args []string) error {
		// init config
		conf := config.GetConfig(cfgFile)

		// init alert bot
		log.Info().Msg("init bot...")
		if conf.AlertConfig.LarkBotId != "" {
			if conf.AlertConfig.Interval != 0 {
				share.SetAlertLimiter(time.Duration(conf.AlertConfig.Interval) * time.Second)
			}
			bot.RegisterNotificator(bot.NewLarkBot(conf.AlertConfig.LarkBotId, ""))
		}

		// init service
		log.Info().Msg("init parser...")
		chainParsers := make(map[string]chain.Parser)
		for chainIndex, info := range conf.EVMEndpoints {
			chainParsers[chainIndex] = chain.NewEvm(info.SupportDebug, info.Endpoint, chainIndex)
		}
		for chainIndex, endpoint := range conf.SOlEndpoints {
			chainParsers[chainIndex] = chain.NewSolana(endpoint, chainIndex)
		}
		for chainIndex, info := range conf.TRONEndpoints {
			chainParsers[chainIndex] = chain.NewTron(info.Endpoint, info.Token, chainIndex)
		}
		chainSvc := service.NewChain(chainParsers)

		// init handler
		log.Info().Msg("init handler...")
		handle := handler.NewHandle(chainSvc)

		gin.SetMode(gin.ReleaseMode)
		g := gin.Default()
		handle.Router(g)
		log.Info().Str("listen", conf.ListenHost).Msg("running...")
		return g.Run(conf.ListenHost)
	},
}

func NewCommand() *cobra.Command {
	serviceCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is /config.yaml)")
	return serviceCmd
}
