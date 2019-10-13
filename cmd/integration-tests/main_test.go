package main

import (
	"flag"
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	integrationTests "github.com/kapustkin/go_calendar/pkg/integration-tests"
)

func TestMain(m *testing.M) {
	flag.Parse()

	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		integrationTests.FeatureContext(s)
	}, godog.Options{
		Output:    colors.Colored(os.Stdout),
		Format:    "pretty",
		Paths:     []string{"features"},
		Randomize: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}