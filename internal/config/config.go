package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"
)

type Config struct {
	SensorUrls     []string      `json:"sensor_urls"`
	SensorTimeoutS string        `json:"sensor_timeout"`
	SensorTimeout  time.Duration `json:"-"`
	IntervalS      string        `json:"interval"`
	Interval       time.Duration `json:"-"`
	Bind           string        `json:"bind"`
	Log            struct {
		Pretty bool   `json:"pretty"`
		Level  string `json:"level"`
	} `json:"log"`
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
	if C.Log.Level == "" {
		C.Log.Level = "info"
	}
	C.Interval, err = time.ParseDuration(C.IntervalS)
	if err != nil {
		return err
	}
	C.SensorTimeout, err = time.ParseDuration(C.SensorTimeoutS)
	return err
}

func init() {
	configPath := flag.String("config", "config.json", "path to config file in json format")
	pretty := flag.Bool("pretty", false, "enable pretty logging")
	level := flag.String(
		"level",
		"",
		"logging level (trace, debug, info, warn, error, fatal, panic, disabled)",
	)
	flag.Parse()
	if err := parseConfig(*configPath); err != nil {
		log.Fatalln("error config:", err)
	}
	if *pretty && !C.Log.Pretty {
		C.Log.Pretty = true
	}
	if *level != "" {
		C.Log.Level = *level
	}
}
