package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/kapustkin/go_calendar/pkg/integration-tests/config"
	"github.com/kapustkin/go_calendar/pkg/service/event-sender/kafka"
)

type NotifyTest struct {
	// config
	config *config.Config
	// kafka
	kafkaConn      *kafka.Kafka
	messages       [][]byte
	messagesMutex  *sync.RWMutex
	stopSignal     chan struct{}
	recievedSignal chan struct{}
	// rest
	responseStatusCode int
	responseBody       []byte
	responseUUID       string
}

func Init(conf *config.Config) *NotifyTest {
	return &NotifyTest{config: conf}
}

func (test *NotifyTest) iSendRequestTo(httpMethod, addr string) (err error) {
	var r *http.Response
	addr = strings.Replace(addr, "{REST_SERVER}", test.config.RestServer, -1)
	switch httpMethod {
	case http.MethodGet:
		//nolint:gosec
		r, err = http.Get(addr)
		if err == nil {
			defer r.Body.Close()
		}
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return err
	}

	test.responseStatusCode = r.StatusCode
	test.responseBody, err = ioutil.ReadAll(r.Body)
	return
}

func (test *NotifyTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *NotifyTest) theResponseShouldMatchText(text string) error {
	if string(test.responseBody) != text {
		return fmt.Errorf("unexpected text: %s != %s", test.responseBody, text)
	}
	return nil
}

func (test *NotifyTest) theResponseShouldContainsText(text string) error {
	if !strings.Contains(string(test.responseBody), text) {
		return fmt.Errorf("unexpected text: %s not contains %s", test.responseBody, text)
	}
	return nil
}

func (test *NotifyTest) theSendRequestToWithData(httpMethod, addr,
	contentType string, data *gherkin.DocString) (err error) {
	var r *http.Response
	addr = strings.Replace(addr, "{REST_SERVER}", test.config.RestServer, -1)
	switch httpMethod {
	case http.MethodPost:
		replacer := strings.NewReplacer("\n", "", "\t", "")
		cleanJSON := replacer.Replace(data.Content)
		//nolint:gosec
		r, err = http.Post(addr, contentType, bytes.NewReader([]byte(cleanJSON)))
		if err == nil {
			defer r.Body.Close()
		}
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return err
	}
	test.responseStatusCode = r.StatusCode
	test.responseBody, err = ioutil.ReadAll(r.Body)
	return
}
