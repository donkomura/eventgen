package eventgen

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"time"
)

type Generator struct {
	Config
	KinesisIterator  func(i int) interface{}
	DynamoDBIterator func(i int) DynamoDBImages
}

func New(c Config) *Generator {
	return &Generator{
		Config: c,
	}
}

func (g *Generator) RegisterKinesis(f func(i int) interface{}) *Generator {
	g.KinesisIterator = f
	return g
}

func (g *Generator) RegisterDynamoDB(f func(i int) DynamoDBImages) *Generator {
	g.DynamoDBIterator = f
	return g
}

func (g *Generator) Kinesis(n int) (*events.KinesisEvent, error) {
	var res []events.KinesisEventRecord
	for i := 0; i < n; i++ {
		b, err := json.Marshal(g.KinesisIterator)
		if err != nil {
			return nil, fmt.Errorf("json marshal: %w", err)
		}

		seqNo := fmt.Sprintf("%057d", i)
		record := events.KinesisRecord{
			ApproximateArrivalTimestamp: events.SecondsEpochTime{Time: time.Now()},
			Data:                        b,
			PartitionKey:                g.PartitionKey,
			SequenceNumber:              seqNo,
		}

		// we ommit `EventVersion`, `EventSourceArn`, and `InvokeIdentityArn`
		res = append(res, events.KinesisEventRecord{
			AwsRegion:      g.Region,
			EventID:        fmt.Sprintf("%s:%s", g.EventIDPrefix, seqNo),
			EventName:      g.EventName,
			EventSource:    g.EventSource,
			EventSourceArn: g.EventSource,
			Kinesis:        record,
		})
	}
	return &events.KinesisEvent{
		Records: res,
	}, nil
}

func (g *Generator) DynamoDB(n int) (*events.DynamoDBEvent, error) {
	var res []events.DynamoDBEventRecord
	for i := 0; i < n; i++ {
		image := g.DynamoDBIterator(i)
		seqNo := fmt.Sprintf("%057d", i)
		record := events.DynamoDBEventRecord{
			AWSRegion:   g.Region,
			EventID:     fmt.Sprintf("%s:%s", g.EventIDPrefix, seqNo),
			EventName:   string(image.EventType),
			EventSource: DynamoDBDefaultEventSource,
		}

		record.Change = events.DynamoDBStreamRecord{
			ApproximateCreationDateTime: events.SecondsEpochTime{Time: time.Now()},
			Keys:                        image.Keys,
			NewImage:                    image.NewImage,
			OldImage:                    image.OldImage,
			SequenceNumber:              seqNo,
			StreamViewType:              image.StreamViewType,
		}
		record.Change.SizeBytes = int64(len(fmt.Sprintf("%+v", record.Change)))
		res = append(res, record)
	}

	return &events.DynamoDBEvent{Records: res}, nil
}
