# eventgen
AWS lambda event generator for Golang tests

# Example

## Kinesis

```go
type MyData struct {
	a int
	b string
}

func main() {
    // register data creation function in each steps 
    generator := eventgen.New(eventgen.DefaultKinesisConfig()).Register(func(i int) interface{} {
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

## DynamoDB stream

```go
generator := eventgen.New(eventgen.DefaultDynamoDBConfig())

	f := func(i int) eventgen.DynamoDBImages {
		return eventgen.DynamoDBImages{
			Keys: map[string]events.DynamoDBAttributeValue{
				"Id": events.NewNumberAttribute("101"),
			},
			NewImage: map[string]events.DynamoDBAttributeValue{
				"Id":      events.NewNumberAttribute("101"),
				"Message": events.NewStringAttribute("New item!"),
			},
			OldImage: map[string]events.DynamoDBAttributeValue{
				"Id":      events.NewNumberAttribute("101"),
				"Message": events.NewStringAttribute("This message has changed"),
			},
            EventType = eventgen.DynamoDBEventTypeModify
			StreamViewType: "NEW_AND_OLD_IMAGES",
		}
	}
	generator.DynamoDBIterator = f

	d, err := generator.DynamoDB(10)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range d.Records {
		fmt.Println(v.EventName)
	}
```
