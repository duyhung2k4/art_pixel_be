package main

import (
	"app/config"
	"app/rabbitmq"
	"app/router"
	"app/socket"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

func runPythonServer() {
	cmd := exec.Command("python3", "python_nodes/index.py")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start Python server: %v", err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Python server exited with error: %v", err)
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		runPythonServer() // Chạy server Python
	}()

	go func() {
		defer wg.Done()
		server := http.Server{
			Addr:           ":" + config.GetAppPort(),
			Handler:        router.AppRouter(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		log.Fatalln(server.ListenAndServe())
	}()

	go func() {
		defer wg.Done()
		socker := http.Server{
			Addr:           ":" + config.GetSocketPort(),
			Handler:        socket.ServerSocker(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		log.Fatalln(socker.ListenAndServe())
	}()

	go func() {
		defer wg.Done()
		rabbitmq.RunRabbitmq()
	}()

	wg.Wait()
}
