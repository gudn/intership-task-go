package worker

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gudn/intership-task-go/pkg/value"
	"github.com/rs/zerolog/log"
)

func fetchValue(req *http.Request) (int32, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	val, err := strconv.ParseUint(string(body), 10, 8)
	if err != nil {
		return 0, err
	}
	return int32(val), err
}

func Start(
	ctx context.Context,
	value *value.Value,
	interval time.Duration,
	url string,
) error {
	atomic.AddInt32(value.BrokenCount, 1)
	lastValue, broken := int32(0), true
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	tick := time.NewTicker(interval)
	defer tick.Stop()
	defer log.Info().Str("url", url).Msg("terminated")
	log.Info().Str("url", url).Msg("initialized")

	update := func() error {
		val, err := fetchValue(req)
		if err != nil {
			if !broken {
				atomic.AddInt32(value.BrokenCount, 1)
				atomic.AddInt32(value.Sum, -lastValue)
				lastValue = 0
				broken = true
			}
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			} else {
				log.Warn().Err(err).Str("url", url).Msg("fetching failed")
			}
		} else {
			atomic.AddInt32(value.Sum, val-lastValue)
			if broken {
				broken = false
				atomic.AddInt32(value.BrokenCount, -1)
			}
			lastValue = val
			log.Info().Str("url", url).Msg("received value")
		}
		return nil
	}

	if err := update(); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			if !broken {
				atomic.AddInt32(value.BrokenCount, 1)
			}
			atomic.AddInt32(value.Sum, -lastValue)
			return ctx.Err()
		case <-tick.C:
			if err := update(); err != nil {
				return err
			}
		}
	}
}
