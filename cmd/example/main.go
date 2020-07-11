package main

import (
	"fmt"
	"github.com/donkomura/eventgen"
	"log"
)

type MyData struct {
	A int
	B string
}

func main() {
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
	fmt.Printf("%+v\n", string(event.Records[0].Kinesis.Data))
}
