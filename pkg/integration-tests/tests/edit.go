package tests

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/google/uuid"
)

func ExecEditTest(s *godog.Suite, test *NotifyTest) {

	s.Step(`^посылаю "([^"]*)" запрос к "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^ожидаю что код ответа будет (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^в ответе будет событие с Message:$`, test.theResponseShouldMatchJSONEvent)
	s.Step(`^посылаю "([^"]*)" запрос к "([^"]*)" c "([^"]*)" и новым содержимым:$`,
		test.theSendRequestToWithDataReplaceUUID)
	s.Step(`^в ответе будет измененное событие Message:$`,
		test.theResponseShouldMatchJSONEvent)
}

type event struct {
	UUID      uuid.UUID
	EventDate string
	Message   string
}

func (test *NotifyTest) theResponseShouldMatchJSONEvent(body *gherkin.DocString) (err error) {
	var expected event
	var actual []event

	// re-encode expected response
	if err = json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return
	}

	// re-encode actual response too
	if err = json.Unmarshal(test.responseBody, &actual); err != nil {
		return
	}

	for _, dt := range actual {
		if dt.Message == expected.Message {
			test.responseUUID = dt.UUID.String()
			return nil
		}
	}

	return fmt.Errorf("event not found")
}

func (test *NotifyTest) theSendRequestToWithDataReplaceUUID(httpMethod, addr,
	contentType string, data *gherkin.DocString) (err error) {

	replacer := strings.NewReplacer("{REPLACE_UUID}", test.responseUUID)
	data.Content = replacer.Replace(data.Content)
	return test.theSendRequestToWithData(httpMethod, addr, contentType, data)
}
