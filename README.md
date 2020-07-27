# eventgen
AWS lambda event generator for Golang tests

# Example

## Kinesis

```go
type MyData struct {
	A int
	B string
}

func main() {
    // register data creation function in each steps 
    generator := eventgen.New(eventgen.DefaultKinesisConfig()).RegisterKinesis(func(i int) interface{} {
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
func main() {
	generator := eventgen.New(eventgen.DefaultDynamoDBConfig()).
		RegisterDynamoDB(func(i int) eventgen.DynamoDBImages {
			e := eventgen.DynamoDBImages{
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
				StreamViewType: "NEW_AND_OLD_IMAGES",
			}
			if i < 3 {
				e.EventType = eventgen.DynamoDBEventTypeInsert
			} else if i < 6 {
				e.EventType = eventgen.DynamoDBEventTypeModify
			} else {
				e.EventType = eventgen.DynamoDBEventTypeMemove
			}

			return e
		})

	d, err := generator.DynamoDB(10)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range d.Records {
		fmt.Println(v.EventName)
	}
}
```
