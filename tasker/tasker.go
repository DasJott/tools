package tasker

import "sync"

// Tasker for running asynchronous tasks
type Tasker struct {
	tasks []func()
}

// Task creates a new Tasker and adds an async task
func Task(f func()) *Tasker {
	t := &Tasker{}
	return t.Task(f)
}

// Task adds an async task
func (t *Tasker) Task(f func()) *Tasker {
	t.tasks = append(t.tasks, f)
	return t
}

// Wait blocks execution while wating for tasks to finish
func (t *Tasker) Wait() {
	if count := len(t.tasks); count > 0 {
		waiter := sync.WaitGroup{}
		waiter.Add(count)
		t.start(&waiter, count)
		waiter.Wait()
	}
}

// Channel returns a channel which is called on finish
func (t *Tasker) Channel() chan int {
	cb := make(chan int)
	go func() {
		t.Wait()
		cb <- len(t.tasks)
	}()
	return cb
}

func (t *Tasker) start(waiter *sync.WaitGroup, count int) {
	for i := 0; i < count; i++ {
		go func() {
			t.tasks[i]()
			waiter.Done()
		}()
	}
}