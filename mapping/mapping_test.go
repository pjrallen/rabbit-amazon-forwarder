package mapping

import (
	"errors"
	"os"
	"testing"

	"github.com/AirHelp/rabbit-amazon-forwarder/common"
	"github.com/AirHelp/rabbit-amazon-forwarder/consumer"
	"github.com/AirHelp/rabbit-amazon-forwarder/forwarder"
	"github.com/AirHelp/rabbit-amazon-forwarder/rabbitmq"
	"github.com/AirHelp/rabbit-amazon-forwarder/sns"
	"github.com/AirHelp/rabbit-amazon-forwarder/sqs"
)

const (
	rabbitType = "rabbit"
	snsType    = "sns"
)

func TestLoad(t *testing.T) {
	os.Setenv(common.MappingFile, "../tests/rabbit_to_sns.json")
	client := New(MockMappingHelper{})
	var consumerForwarderMap map[consumer.Client]forwarder.Client
	var err error
	if consumerForwarderMap, err = client.Load(); err != nil {
		t.Errorf("could not load mapping and start mocked rabbit->sns pair: %s", err.Error())
	}
	if len(consumerForwarderMap) != 1 {
		t.Errorf("wrong consumerForwarderMap size, expected 1, got %d", len(consumerForwarderMap))
	}
}

func TestLoadFile(t *testing.T) {
	os.Setenv(common.MappingFile, "../tests/rabbit_to_sns.json")
	client := New()
	data, err := client.loadFile()
	if err != nil {
		t.Errorf("could not load file: %s", err.Error())
	}
	if len(data) < 1 {
		t.Errorf("could not load file: empty steam found")
	}
}

func TestCreateConsumer(t *testing.T) {
	client := New()
	consumerName := "test-rabbit"
	item := common.Item{Type: "RabbitMQ",
		Name:          consumerName,
		ConnectionURL: "url",
		ExchangeName:  "topic",
		QueueName:     "test-queue",
		RoutingKey:    "#"}
	consumer := client.helper.createConsumer(item)
	if consumer.Name() != consumerName {
		t.Errorf("wrong consumer name, expected %s, found %s", consumerName, consumer.Name())
	}
}

func TestCreateForwarderSNS(t *testing.T) {
	client := New(MockMappingHelper{})
	forwarderName := "test-sns"
	item := common.Item{Type: "SNS",
		Name:          forwarderName,
		ConnectionURL: "",
		ExchangeName:  "topic",
		QueueName:     "",
		RoutingKey:    "#"}
	forwarder := client.helper.createForwarder(item)
	if forwarder.Name() != forwarderName {
		t.Errorf("wrong forwarder name, expected %s, found %s", forwarderName, forwarder.Name())
	}
}

func TestCreateForwarderSQS(t *testing.T) {
	client := New(MockMappingHelper{})
	forwarderName := "test-sqs"
	item := common.Item{Type: "SQS",
		Name:          forwarderName,
		ConnectionURL: "",
		ExchangeName:  "",
		QueueName:     "test-queue",
		RoutingKey:    "#"}
	forwarder := client.helper.createForwarder(item)
	if forwarder.Name() != forwarderName {
		t.Errorf("wrong forwarder name, expected %s, found %s", forwarderName, forwarder.Name())
	}
}

// helpers
type MockMappingHelper struct{}

type MockRabbitConsumer struct{}

type MockSNSForwarder struct {
	name string
}

type MockSQSForwarder struct {
	name string
}

type ErrorForwarder struct{}

func (h MockMappingHelper) createConsumer(item common.Item) consumer.Client {
	if item.Type != rabbitmq.Type {
		return nil
	}
	return MockRabbitConsumer{}
}
func (h MockMappingHelper) createForwarder(item common.Item) forwarder.Client {
	switch item.Type {
	case sns.Type:
		return MockSNSForwarder{item.Name}
	case sqs.Type:
		return MockSQSForwarder{item.Name}
	}
	return ErrorForwarder{}
}

func (c MockRabbitConsumer) Name() string {
	return rabbitType
}

func (c MockRabbitConsumer) Start(client forwarder.Client, check chan bool, stop chan bool) error {
	return nil
}

func (f MockSNSForwarder) Name() string {
	return f.name
}

func (f MockSNSForwarder) Push(message string) error {
	return nil
}

func (f MockSQSForwarder) Name() string {
	return f.name
}

func (f MockSQSForwarder) Push(message string) error {
	return nil
}

func (f ErrorForwarder) Name() string {
	return "error-forwarder"
}

func (f ErrorForwarder) Push(message string) error {
	return errors.New("Wrong forwader created")
}
