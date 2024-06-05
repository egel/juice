package main

import (
	"github.com/egel/juice/pkg/juice"
	"github.com/egel/juice/pkg/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	output := logger.NewZerologConsoleWriter()
	log.Logger = log.Output(output).With().Caller().Logger()

	juice.Execute()
}
