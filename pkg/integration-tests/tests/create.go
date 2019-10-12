package tests

import (
	"github.com/DATA-DOG/godog"
)

func ExecCreateTest(s *godog.Suite, test *notifyTest) {
	s.Step(`^посылаю "([^"]*)" запрос к "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^ожидаю что код ответа будет (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^тело ответа будет равно "([^"]*)"$`, test.theResponseShouldMatchText)

	s.Step(`^посылаю "([^"]*)" запрос к "([^"]*)" c "([^"]*)" содержимым:$`, test.theSendRequestToWithData)
	s.Step(`^в ответе будет событие с Message "([^"]*)"$`, test.theResponseShouldContainsText)
}
