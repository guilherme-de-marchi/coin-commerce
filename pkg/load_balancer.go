package pkg

import (
	"context"
	"sync"
)

type WorkerFunc func(context.Context, *sync.WaitGroup, LoadBalancerRequestChannel) error

type LoadBalancerRequest struct {
	Exchange,
	Target string
	Data         any
	ResponseChan chan []byte
}

type LoadBalancerRequestChannel chan LoadBalancerRequest

type LoadBalancer struct {
	RequestChannel chan LoadBalancerRequest
}

func NewLoadBalancer(reqChan LoadBalancerRequestChannel) *LoadBalancer {
	return &LoadBalancer{RequestChannel: reqChan}
}

func (l *LoadBalancer) Start(ctx context.Context, workers int, wf WorkerFunc) error {
	var channels []LoadBalancerRequestChannel

	wg := new(sync.WaitGroup)
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		c := make(LoadBalancerRequestChannel, 1)
		go wf(ctx, wg, c)
		channels = append(channels, c)
	}

	stopIteratorChan := make(chan bool, 1)
	go loadBalancerIterator(workers, channels, l.RequestChannel, stopIteratorChan)

	wg.Wait()
	stopIteratorChan <- true

	return nil
}

func loadBalancerIterator(workers int, channels []LoadBalancerRequestChannel, reqChan LoadBalancerRequestChannel, stopChan chan bool) {
	running := true
	go func() {
		var i int
		for running {
			if i == workers {
				i = 0
			}

			req := <-reqChan
			channels[i] <- req
		}
	}()

	<-stopChan
	running = false
	return
}
