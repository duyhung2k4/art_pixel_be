package rabbitmq

import (
	"log"
	"sync"
)

func RunRabbitmq() {
	authQueue := NewQueueAuth()
	eventQueue := NewQueueEvent()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		authQueue.InitQueueSendFileAuth()
		wg.Done()
	}()
	go func() {
		authQueue.InitQueueAuthFace()
		wg.Done()
	}()
	go func() {
		eventQueue.InitQueueDrawPixel()
		wg.Done()
	}()

	log.Println("run rabbitmq successfully!")
	wg.Wait()

}
