package tests

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/kapustkin/go_calendar/pkg/service/event-sender/config"
	"github.com/kapustkin/go_calendar/pkg/service/event-sender/kafka"
)

func ExecCreateTest(s *godog.Suite, test *NotifyTest) {
	s.BeforeScenario(test.startKafkaConsuming)
	s.Step(`^посылаю "([^"]*)" запрос к "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^ожидаю что код ответа будет (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^тело ответа будет равно "([^"]*)"$`, test.theResponseShouldMatchText)

	s.Step(`^посылаю "([^"]*)" запрос к "([^"]*)" c "([^"]*)" содержимым:$`, test.theSendRequestToWithData)
	s.Step(`^в ответе будет событие с Message "([^"]*)"$`, test.theResponseShouldContainsText)
	s.Step(`^дождаться оповещения о событии с сообщением "([^"]*)"$`, test.iReceiveEventWithText)
}

func (test *NotifyTest) startKafkaConsuming(interface{}) {
	// read env config
	conf := &config.Config{
		KafkaConnection: test.config.KafkaConnection,
		KafkaTopic:      test.config.KafkaTopic,
		//KafkaPartition:  1,
	}
	// init kafka
	kafkaConn, err := kafka.Init(conf)
	panicOnErr(err)
	test.kafkaConn = kafkaConn

	test.messages = make([][]byte, 0)
	test.messagesMutex = new(sync.RWMutex)
	test.stopSignal = make(chan struct{})
	test.recievedSignal = make(chan struct{}, 1)
	go func(stop <-chan struct{}) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		for {
			select {
			case <-stop:
				return
			default:
				{
					message, err := kafkaConn.GetMessage(ctx)
					panicOnErr(err)
					test.messagesMutex.Lock()
					test.messages = append(test.messages, message)
					test.messagesMutex.Unlock()
					test.recievedSignal <- struct{}{}
				}
			}
		}
	}(test.stopSignal)
}

func (test *NotifyTest) iReceiveEventWithText(text string) error {
	<-test.recievedSignal
	test.messagesMutex.RLock()
	defer test.messagesMutex.RUnlock()

	//test.stopSignal <- struct{}{}
	//panicOnErr(test.kafkaConn.Close())

	for _, msg := range test.messages {
		if strings.Contains(string(msg), text) {
			return nil
		}
	}

	return fmt.Errorf("event with text '%s' was not found in %s", text, test.messages)
}

func panicOnErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
