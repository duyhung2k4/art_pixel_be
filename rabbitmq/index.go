package rabbitmq

import (
	"log"
	"sync"
)

func RunRabbitmq() {
	authQueue := NewQueueAuth()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		authQueue.InitQueueSendFileAuth()
		wg.Done()
	}()
	go func() {
		authQueue.InitQueueAuthFace()
		wg.Done()
	}()

	log.Println("run rabbitmq successfully!")
	wg.Wait()

}
