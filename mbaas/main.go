package main

import (
	"github.com/joeshaw/envdecode"
	"github.com/rs/zerolog/log"

	"github.com/ian-ross/morse-blinkies/mbaas/chassis"
	"github.com/ian-ross/morse-blinkies/mbaas/server"
)

func main() {
	cfg := server.Config{}
	err := envdecode.StrictDecode(&cfg)
	if err != nil {
		log.Fatal().Err(err).
			Msg("failed to process environment variables")
	}
	chassis.LogSetup(cfg.DevMode)
	serv := server.NewServer(&cfg)
	serv.Serve()
}
