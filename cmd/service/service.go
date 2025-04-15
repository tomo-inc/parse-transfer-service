package service

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tomo.inc/parse-transfer-service/app/handler"
	"github.com/tomo.inc/parse-transfer-service/app/service"
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
		if conf.AlertConfig.LarkBotId != "" {
			if conf.AlertConfig.Interval != 0 {
				share.SetAlertLimiter(time.Duration(conf.AlertConfig.Interval) * time.Second)
			}
			bot.RegisterNotificator(bot.NewLarkBot(conf.AlertConfig.LarkBotId, ""))
		}

		// init service
		evmChainConfig := make(map[string]service.EvmConfig)
		for chainIndex, info := range conf.EVMEndpoints {
			evmChainConfig[chainIndex] = service.EvmConfig{
				Endpoint:     info.Endpoint,
				SupportDebug: info.SupportDebug,
			}
		}
		chainSvc := service.NewChain(evmChainConfig, conf.SOlEndpoints, conf.TRONEndpoints)

		// init handler
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
