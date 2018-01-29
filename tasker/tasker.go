package tasker

import (
	"sync"
)

type Task func()
type TaskReturn func() interface{}

// Tasker for running asynchronous tasks
type Tasker struct {
	tasks   []interface{}
	mutex   sync.Mutex
	results []interface{}
	later   chan int
}

// Add creates a new Tasker and adds one or more async tasks. Takes func()  or  func() interface{}
func Add(f ...interface{}) *Tasker {
	t := &Tasker{}
	return t.Add(f...)
}

// Add adds one or more async tasks. Takes func()  or  func() interface{}
func (t *Tasker) Add(f ...interface{}) *Tasker {
	t.tasks = append(t.tasks, f...)
	return t
}

// Wait blocks execution while wating for tasks to finish
func (t *Tasker) Wait() *Tasker {
	if count := len(t.tasks); count > 0 {
		waiter := sync.WaitGroup{}
		waiter.Add(count)
		t.start(&waiter, count)
		waiter.Wait()
	}
	return t
}

// Later returns immediately. Call Now() if you want to continue.
func (t *Tasker) Later() *Tasker {
	t.later = make(chan int)
	go func() {
		t.Wait()
		t.later <- len(t.tasks)
	}()
	return t
}

// Now blocks and returns after call to Later().
func (t *Tasker) Now() *Tasker {
	<-t.later
	return t
}

// Get gets all results collected by using TaskReturn functions
func (t *Tasker) Get() []interface{} {
	return t.results
}

func (t *Tasker) start(waiter *sync.WaitGroup, count int) {
	for i := 0; i < count; i++ {
		if task, ok := t.tasks[i].(func()); ok {
			go func() {
				task()
				waiter.Done()
			}()
		} else if task, ok := t.tasks[i].(func() interface{}); ok {
			go func() {
				t.addReturn(task())
				waiter.Done()
			}()
		} else {
			panic("invalid function format")
		}
	}
}

func (t *Tasker) addReturn(ret interface{}) {
	t.mutex.Lock()
	t.results = append(t.results, ret)
	t.mutex.Unlock()
}
