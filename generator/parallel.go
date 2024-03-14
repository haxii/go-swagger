package generator

import (
	"fmt"
	"sync"
)

type ParallelTaskError map[int]error

func (p ParallelTaskError) Error() string {
	msg := ""
	for key, err := range p {
		msg += fmt.Sprintf("%d: %s;", key, err)
	}
	return msg
}

type Task func(index int) error

func ParallelExecute(size, thread int, task Task) error {
	if size <= 0 || task == nil || thread <= 0 {
		return nil
	}
	if thread > size {
		thread = size
	}
	errs := make(ParallelTaskError)
	errMu := sync.Mutex{}
	keyChan := make(chan int, thread)

	go func() {
		for key := range size {
			keyChan <- key
		}
		close(keyChan)
	}()

	wg := sync.WaitGroup{}
	wg.Add(thread)
	for i := 0; i < thread; i++ {
		go func(t int) {
			for key := range keyChan {
				if err := task(key); err != nil {
					errMu.Lock()
					errs[key] = err
					errMu.Unlock()
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	if len(errs) == 0 {
		return nil
	}
	return errs
}
