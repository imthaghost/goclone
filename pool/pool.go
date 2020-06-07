package pool

import (
	"sync"
)

// Pool ...
type Pool struct {
	numGo    int
	messages chan interface{}
	function func(interface{})
}

// New ...
func New(numGoroutine int, function func(interface{})) *Pool {
	return &Pool{
		numGo:    numGoroutine,
		messages: make(chan interface{}),
		function: function,
	}
}

// Push ...
func (c *Pool) Push(data interface{}) {
	c.messages <- data
}

// CloseQueue ...
func (c *Pool) CloseQueue() {
	close(c.messages)
}

// Run ...
func (c *Pool) Run() {
	var wg sync.WaitGroup

	wg.Add(c.numGo)

	for i := 0; i < c.numGo; i++ {
		go func() {
			for v := range c.messages {
				c.function(v)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
