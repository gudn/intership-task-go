package intership_task_go

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gudn/intership-task-go/internal/config"
	"github.com/gudn/intership-task-go/pkg/value"
	"github.com/gudn/intership-task-go/pkg/worker"

	"github.com/rs/zerolog/log"
)

const (
	BrokenHeader = "X-Broken"
)

var val *value.Value

func acceptsJson(r *http.Request) bool {
	value := r.Header.Get("Accept")
	return strings.Contains(value, "*/*") ||
		strings.Contains(value, "application/json")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	avg, broken := val.Average()
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
		log.Error().Err(err).Msg("handler error")
	}
}

func Setup(ctx context.Context) {
	n := len(config.C.SensorUrls)
	if n < 1 {
		log.Panic().Msg("no configured sensors")
	}
	val = value.New(n)
	for _, url := range config.C.SensorUrls {
		go worker.Start(ctx, val, config.C.Interval, url)
	}
}
