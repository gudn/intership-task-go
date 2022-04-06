package worker

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	intership_task_go "github.com/gudn/intership-task-go"
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
	value *intership_task_go.Value,
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
	for {
		select {
		case <-ctx.Done():
			if !broken {
				atomic.AddInt32(value.BrokenCount, 1)
			}
			atomic.AddInt32(value.Sum, -lastValue)
			return ctx.Err()
		case <-tick.C:
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
					log.Printf("error worker(url = %q): %v", url, err)
				}
			} else {
				atomic.AddInt32(value.Sum, val-lastValue)
				if broken {
					broken = false
					atomic.AddInt32(value.BrokenCount, -1)
				}
				lastValue = val
				log.Printf("worker(url = %q) received value", url)
			}
		}
	}
}
