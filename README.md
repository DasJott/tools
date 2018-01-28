#tools

## tasker
one or more tasks running concurrently.
### example
```go
tasker.Task(func() {
	// whatever yiu want to do
}).Task(func() {
	// something else
}).Wait()
// execution blocked until all tasks finish
```
OR
```go
fin := tasker.Task(func() {
	// whatever yiu want to do
}).Task(func() {
	// something else
}).Channel()
// channel 'fin' can be asked whenever you want
```

