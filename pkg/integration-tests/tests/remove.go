package tests

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

func ExecRemoveTest(s *godog.Suite, test *NotifyTest) {
	s.Step(`^в ответе не будет события Message:$`, test.theResponseShouldEventNotFound)
}

func (test *NotifyTest) theResponseShouldEventNotFound(body *gherkin.DocString) (err error) {
	var expected event
	var actual []event

	replacer := strings.NewReplacer("{REPLACE_UUID}", test.responseUUID)
	body.Content = replacer.Replace(body.Content)

	// re-encode expected response
	if err = json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return
	}

	// re-encode actual response too
	if err = json.Unmarshal(test.responseBody, &actual); err != nil {
		return
	}

	for _, dt := range actual {
		if dt.UUID == expected.UUID {
			test.responseUUID = dt.UUID.String()
			return fmt.Errorf("event found")
		}
	}
	return nil
}
