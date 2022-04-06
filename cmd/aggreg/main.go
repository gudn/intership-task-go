package main

import (
	"context"
	"net"
	"net/http"

	intership_task_go "github.com/gudn/intership-task-go"
	"github.com/gudn/intership-task-go/internal/config"
	_ "github.com/gudn/intership-task-go/internal/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	intership_task_go.Setup(ctx)
	http.HandleFunc("/", intership_task_go.Handler)
	listen, err := net.Listen("tcp", config.C.Bind)
	if err != nil {
		log.Fatal().Err(err).Msg("error listening")
	}
	log.Info().Str("bind", listen.Addr().String()).Msg("started listening")
	if err := http.Serve(listen, nil); err != nil {
		log.Fatal().Err(err).Msg("listening failed")
	}
}
