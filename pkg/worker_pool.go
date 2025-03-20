package pkg

import (
	"fmt"
	"sync"
)

type Pool struct {
	capacity         int                             // number of go routines
	messageQueueChan chan interface{}                // queue for contains message
	wg               sync.WaitGroup                  // used to graceful shutdown
	processFunc      func(message interface{}) error // function to process message
}

func NewPool(capacity, messageQueueSize int, processFunc func(message interface{}) error) *Pool {
	return &Pool{
		capacity:         capacity,
		messageQueueChan: make(chan interface{}, messageQueueSize),
		processFunc:      processFunc,
	}
}

func (p *Pool) PushMessage(message interface{}) {
	p.messageQueueChan <- message
}

func (p *Pool) Start() {
	for i := 0; i < p.capacity; i++ {
		p.wg.Add(1)
		go p.handleMessage(i, p.messageQueueChan)
	}
}

func (p *Pool) handleMessage(id int, msgChan chan interface{}) {
	defer p.wg.Done()

	for message := range msgChan {
		fmt.Printf("worker %v is processing message %v\n", id, message)
		err := p.processFunc(message)

		if err != nil {
			fmt.Printf("worker %v failed to process message %v with error: %v\n", id, message, err)
		}
	}
}

// graceful shutdown
func (p *Pool) GracefulShutdown() {
	fmt.Println("Closing worker pool...")
	close(p.messageQueueChan)
	p.wg.Wait()
	fmt.Println("Worker pool gracefully shutdown!")
}
