package eventgen

import "github.com/aws/aws-lambda-go/events"

const (
	DefaultRegion        = "us-east-2"
	DefaultEventIDPrefix = "eventgen:"

	KinesisDefaultEventSource  = "eventgen:kinesis"
	KinesisDefaultEventName    = KinesisDefaultEventSource + ":record"
	KinesisDefaultPartitionKey = "1"

	DynamoDBDefaultEventSource = "eventgen:dynamodb"
)

type DynamoDBEventType string

var (
	DynamoDBEventTypeInsert = DynamoDBEventType("INSERT")
	DynamoDBEventTypeModify = DynamoDBEventType("MODIFY")
	DynamoDBEventTypeMemove = DynamoDBEventType("REMOVE")
)

type DynamoDBImages struct {
	Keys           map[string]events.DynamoDBAttributeValue
	NewImage       map[string]events.DynamoDBAttributeValue
	OldImage       map[string]events.DynamoDBAttributeValue
	StreamViewType string
	EventType      DynamoDBEventType
}

type Config struct {
	Region        string
	EventIDPrefix string
	EventName     string
	EventSource   string
	PartitionKey  string
}

func DefaultKinesisConfig() Config {
	return Config{
		Region:        DefaultRegion,
		EventIDPrefix: DefaultEventIDPrefix,
		EventName:     KinesisDefaultEventName,
		EventSource:   KinesisDefaultEventSource,
		PartitionKey:  KinesisDefaultPartitionKey,
	}
}

func DefaultDynamoDBConfig() Config {
	return Config{
		Region:        DefaultRegion,
		EventIDPrefix: DefaultEventIDPrefix,
		EventSource:   DynamoDBDefaultEventSource,
	}
}
