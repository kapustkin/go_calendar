package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/kapustkin/go_calendar/pkg/service/event-sender/kafka"
)

type notifyTest struct {
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

func Init() *notifyTest {
	return &notifyTest{}
}

func (test *notifyTest) iSendRequestTo(httpMethod, addr string) (err error) {
	var r *http.Response

	switch httpMethod {
	case http.MethodGet:
		r, err = http.Get(addr)
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return
	}
	test.responseStatusCode = r.StatusCode
	test.responseBody, err = ioutil.ReadAll(r.Body)
	return
}

func (test *notifyTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *notifyTest) theResponseShouldMatchText(text string) error {
	if string(test.responseBody) != text {
		return fmt.Errorf("unexpected text: %s != %s", test.responseBody, text)
	}
	return nil
}

func (test *notifyTest) theResponseShouldContainsText(text string) error {
	if !strings.Contains(string(test.responseBody), text) {
		return fmt.Errorf("unexpected text: %s not contains %s", test.responseBody, text)
	}
	return nil
}

func (test *notifyTest) theSendRequestToWithData(httpMethod, addr, contentType string, data *gherkin.DocString) (err error) {
	var r *http.Response

	switch httpMethod {
	case http.MethodPost:
		replacer := strings.NewReplacer("\n", "", "\t", "")
		cleanJson := replacer.Replace(data.Content)
		r, err = http.Post(addr, contentType, bytes.NewReader([]byte(cleanJson)))
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return
	}
	test.responseStatusCode = r.StatusCode
	test.responseBody, err = ioutil.ReadAll(r.Body)
	return
}

func (test *notifyTest) theResponseShouldMatchJSON(body *gherkin.DocString) (err error) {
	var expected, actual interface{}

	// re-encode expected response
	if err = json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return
	}

	// re-encode actual response too
	if err = json.Unmarshal(test.responseBody, &actual); err != nil {
		return
	}

	// the matching may be adapted per different requirements.
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}
	return nil
}
