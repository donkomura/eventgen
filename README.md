# eventgen
AWS lambda event generator for Golang tests

# Example

```go
type MyData struct {
	a int
	b string
}

func main() {
	generator := eventgen.New(MyData{}).Register(func(i int) MyData {
		// generate data for each step
	})

	kinesisEvent := generator.Kinesis(100) // create 100 records in a event

	// use evnets here
}
```
