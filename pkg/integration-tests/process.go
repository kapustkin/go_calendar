package integrationtests

import (
	"log"

	"github.com/DATA-DOG/godog"
	"github.com/kapustkin/go_calendar/pkg/integration-tests/tests"
)

func FeatureContext(s *godog.Suite) {
	// загрузка конфига
	test := tests.Init()
	// выход из сценария, если он завершился с ошибкой
	s.AfterScenario(func(data interface{}, err error) {
		if err != nil {
			log.Fatalf("%v", err)
		}
	})
	tests.ExecCreateTest(s, test)
	tests.ExecEditTest(s, test)
	tests.ExecRemoveTest(s, test)
}
