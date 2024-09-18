package rabbitmq

import (
	"log"
	"sync"
)

func RunRabbitmq() {
	authQueue := NewQueueAuth()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		authQueue.InitQueueSendFileAuth()
		wg.Done()
	}()

	log.Println("run rabbitmq successfully!")
	wg.Wait()

}
