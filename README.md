# eventgen
AWS lambda event generator for Golang tests

# Example

```go
type MyData struct {
	a int
	b string
}

func main() {
    // register data creation function in each steps 
    generator := eventgen.New(MyData{}).Register(func(i int) interface{} {
		return MyData{
			A: i,
			B: fmt.Sprintf("%dth", i),
		}
	})
	event, err := generator.Kinesis(10)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range event.Records {
		fmt.Printf("%+v\n", v.Kinesis)
	}
}
```
