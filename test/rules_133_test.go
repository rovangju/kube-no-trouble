package test

import (
	"testing"
)

func TestRego133(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"Endpoints", []string{"../fixtures/endpoints-v1.yaml"}, []string{"Endpoints"}},
	}

	testResourcesUsingFixtures(t, testCases)
}
