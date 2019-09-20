#tools

## global
things regarding globalization
### contents
- **ZoneMap**<br>
  Map from timezone names to its offset on UTC (`map[string]time.Duration`)


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
	// whatever you want to do
}).Add(func() {
	// something else
}).Later()
// Later returns immediately and can be continued by calling t.Now()
```

