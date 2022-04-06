package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	intership_task_go "github.com/gudn/intership-task-go"
	"github.com/gudn/intership-task-go/internal/config"
	"github.com/gudn/intership-task-go/pkg/worker"
)

const (
	BrokenHeader = "X-Broken"
)

var value *intership_task_go.Value

func acceptsJson(r *http.Request) bool {
	value := r.Header.Get("Accept")
	return strings.Contains(value, "*/*") ||
		strings.Contains(value, "application/json")
}

func handler(w http.ResponseWriter, r *http.Request) {
	avg, broken := value.Average()
	if broken {
		w.Header().Add(BrokenHeader, "true")
	}
	var err error
	if acceptsJson(r) {
		data := map[string]interface{}{
			"broken":  broken,
			"average": avg,
		}
		var encoded []byte
		encoded, err = json.Marshal(data)
		if err == nil {
			_, err = w.Write(encoded)
		}
	} else {
		_, err = fmt.Fprint(w, avg)
	}
	if err != nil {
		log.Print("error handler:", err)
	}
}

func main() {
	n := len(config.C.SensorUrls)
	if n < 1 {
		log.Fatalln("error config: sensors count is too few")
	}
	value = intership_task_go.NewValue(n)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, url := range config.C.SensorUrls {
		go worker.Start(ctx, value, config.C.Interval, url)
	}
	http.HandleFunc("/", handler)
	listen, err := net.Listen("tcp", config.C.Bind)
	if err != nil {
		log.Fatalln("error listen:", err)
	}
	log.Print("listening on ", listen.Addr().String())
	log.Fatalln(http.Serve(listen, nil))
}
