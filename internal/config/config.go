package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"
)

type Config struct {
	SensorUrls []string      `json:"sensor_urls"`
	IntervalS  string        `json:"interval"`
	Interval   time.Duration `json:"-"`
	Bind       string        `json:"bind"`
}

var C Config

func parseConfig(fname string) error {
	contents, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	err = json.Unmarshal(contents, &C)
	if err != nil {
		return err
	}
	C.Interval, err = time.ParseDuration(C.IntervalS)
	return err
}

func init() {
	configPath := flag.String("config", "config.json", "path to config file in json format")
	flag.Parse()
	if err := parseConfig(*configPath); err != nil {
		log.Fatalln("error config:", err)
	}
}
