package eventgen

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

type Generator struct {
	Struct interface{}
	Iterate func(i int) interface{}
}

func New(i interface{}) *Generator {
	return &Generator{Struct:  i}
}

func (g *Generator) Register(f func(i int) interface{}) *Generator {
	g.Iterate = f
	return g
}

func (g *Generator) Kinesis(n int) (*events.KinesisEvent, error) {
	var res []events.KinesisEventRecord
	for i := 0; i < n; i++ {
		var record events.KinesisRecord
		b, err := json.Marshal(g.Iterate(i))
		if err != nil {
			return nil, fmt.Errorf("json marshal: %w", err)
		}
		record.Data = b
		res = append(res, events.KinesisEventRecord{
			Kinesis: record,
		})
	}
	return &events.KinesisEvent{
		Records: res,
	}, nil
}
