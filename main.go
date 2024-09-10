package main

import (
	"app/router"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		server := http.Server{
			Addr:           ":" + "10000",
			Handler:        router.AppRouter(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		log.Fatalln(server.ListenAndServe())
		wg.Done()
	}()

	wg.Wait()
}
