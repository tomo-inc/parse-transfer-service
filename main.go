package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/tomo.inc/parse-transfer-service/cmd/service"
	"github.com/tomo.inc/parse-transfer-service/cmd/version"
	"log"
	"os"
)

func main() {
	cmd := &cobra.Command{
		Use:   "parse-transfer-service",
		Short: "parse transaction to get transfer",
	}

	cmd.AddCommand(
		service.NewCommand(),
		version.NewCommand(),
	)

	setLogLevel()

	if err := cmd.ExecuteContext(cmd.Context()); err != nil {
		log.Fatal(fmt.Sprintf("excute failed: %v", err))
	}
}

func setLogLevel() {
	switch os.Getenv("LOG_LEVEL") {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}
