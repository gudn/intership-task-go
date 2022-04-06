package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	flag.Parse()
}

var val uint8

func main() {
	rand.Seed(time.Now().Unix())

	go func() {
		for {
			val = uint8(rand.Intn(255))
			log.Printf("value changed, new value is %v", val)
			time.Sleep(time.Second * 30)
		}
	}()

	err := http.ListenAndServe(":8080", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if rand.Intn(20) == 0 {
			time.Sleep(time.Second * time.Duration(rand.Intn(10)))
		}

		_, err := fmt.Fprintf(writer, "%d", val)
		if err != nil {
			log.Print("error", err)
		}
	}))

	if err != nil {
		log.Fatal(err)
	}
}
