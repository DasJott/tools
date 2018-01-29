#tools

## tasker
one or more tasks running concurrently.
### example
```go
tasker.Add(
	func() {
		// whatever yiu want to do
	},
	func() {
		// something else
	},
).Wait()
// execution blocked until all tasks finish
```
OR
```go
t := tasker.Add(func() {
	// whatever yiu want to do
}).Add(func() {
	// something else
}).Later()
// Later returns immediately and can be continued by calling t.Now()
```
