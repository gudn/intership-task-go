package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	intership_task_go "github.com/gudn/intership-task-go"
	"github.com/gudn/intership-task-go/pkg/worker"
)

const (
	BrokenHeader = "X-Broken"
)

type Config struct {
	SensorUrls []string      `json:"sensor_urls"`
	IntervalS  string        `json:"interval"`
	Interval   time.Duration `json:"-"`
	Bind       string        `json:"bind"`
}

var config Config
var value *intership_task_go.Value

func parseConfig(fname string) error {
	contents, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	err = json.Unmarshal(contents, &config)
	if err != nil {
		return err
	}
	config.Interval, err = time.ParseDuration(config.IntervalS)
	return err
}

func acceptsJson(r *http.Request) bool {
	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	avg, broken := value.Average()
	if broken {
		w.Header().Add(BrokenHeader, "true")
	}
	var err error
	if acceptsJson(r) {
		data := map[string]interface{}{
			"broken": broken,
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
	configPath := flag.String("config", "config.json", "path to config file in json format")
	flag.Parse()
	if err := parseConfig(*configPath); err != nil {
		log.Fatalln("error config:", err)
	}
	n := len(config.SensorUrls)
	if n < 1 {
		log.Fatalln("error config: sensors count is too few")
	}
	value = intership_task_go.NewValue(n)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, url := range config.SensorUrls {
		go worker.Start(ctx, value, config.Interval, url)
	}
	http.HandleFunc("/", handler)
	listen, err := net.Listen("tcp", config.Bind)
	if err != nil {
		log.Fatalln("error listen:", err)
	}
	log.Print("listening on ", listen.Addr().String())
	log.Fatalln(http.Serve(listen, nil))
}
