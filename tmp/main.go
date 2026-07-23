package main

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/errgroup"
)

type Count struct {
	value int
	mu    sync.Mutex
}

func (c *Count) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.value++
}

type SafeCount struct {
	value int
	mu    sync.RWMutex
}

func (c *SafeCount) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.value++
}

func (c *SafeCount) Get() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.value
}

func main() {
	/*
		// Mutex (RW)
		counter := new(SafeCount)
		var wg sync.WaitGroup

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				counter.Inc()
			}()
		}
		wg.Wait()

		fmt.Println("Итоговое количество", counter.Get())*/
	/*	q := NewQueue()
		go q.Pop()
		go q.Push(1)*/

	/*	var config atomic.Value

		config.Store(2)
		cfg := config.Load()*/

}

func fetchAll(urls []string) error {
	//gr := new(errgroup.Group)
	g, ctx := errgroup.WithContext(context.Background())
	for _, url := range urls {
		g.Go(func() error {
			req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			return nil
		})
	}
	return g.Wait()
}

//// sync.Once
//var (
//	once     sync.Once
//	instance *Myservice
//)

//	func GetService() *Myservice {
//		once.Do(func() {
//			instance = &Myservice{}
//		})
//		return instance
//	}
//
// sync.Cond
type Queue struct {
	items []int
	cond  *sync.Cond
}

func NewQueue() *Queue {
	return &Queue{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (q *Queue) Push(item int) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.items = append(q.items, item)
	q.cond.Signal()
}

func (q *Queue) Pop() int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for len(q.items) == 0 {
		q.cond.Wait()
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

var pool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func process(data []byte) {
	buf := pool.Get().(*bytes.Buffer)
	defer pool.Put(buf)
	buf.Reset()
	buf.Write(data)
	// работа с буфером
}

///

var counter int32

func Inc() {
	atomic.AddInt32(&counter, 1)
}

func Get() int32 {
	return atomic.LoadInt32(&counter)
}

//
