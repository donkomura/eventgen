package eventgen

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"time"
)

const (
	DefaultEventSource   = "eventgen:kinesis"
	DefaultEventName     = DefaultEventSource + ":record"
	DefaultEventIDPrefix = "eventgen:"
	DefaultRegion        = "us-east-2"
	DefaultPartitionKey  = "1"
)

type Config struct {
	Region        string
	EventIDPrefix string
	EventName     string
	EventSource   string
	PartitionKey  string
}

type Generator struct {
	Config
	Struct  interface{}
	Iterate func(i int) interface{}
}

func New(i interface{}) *Generator {
	return &Generator{
		Config: Config{
			Region:        DefaultRegion,
			EventIDPrefix: DefaultEventIDPrefix,
			EventName:     DefaultEventName,
			EventSource:   DefaultEventSource,
			PartitionKey:  DefaultPartitionKey,
		},
		Struct: i,
	}
}

// optional
func (g *Generator) Setup(c Config) *Generator {
	g.Config = c
	return g
}

func (g *Generator) Register(f func(i int) interface{}) *Generator {
	g.Iterate = f
	return g
}

func (g *Generator) Kinesis(n int) (*events.KinesisEvent, error) {
	var res []events.KinesisEventRecord
	for i := 0; i < n; i++ {
		b, err := json.Marshal(g.Iterate(i))
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
