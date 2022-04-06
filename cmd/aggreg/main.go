package main

import (
	"context"
	"log"
	"net"
	"net/http"

	intership_task_go "github.com/gudn/intership-task-go"
	"github.com/gudn/intership-task-go/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	intership_task_go.Setup(ctx)
	http.HandleFunc("/", intership_task_go.Handler)
	listen, err := net.Listen("tcp", config.C.Bind)
	if err != nil {
		log.Fatalln("error listen:", err)
	}
	log.Print("listening on ", listen.Addr().String())
	log.Fatalln(http.Serve(listen, nil))
}
