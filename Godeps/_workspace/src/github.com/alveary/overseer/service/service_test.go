package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
}

// func (suite *TestSuite) SetupTest() {
// }

func (suite *ServiceTestSuite) TestFailIncrement() {
	testService := Service{"name", "uri", "uri", 0}

	assert.Equal(suite.T(), testService.Fails, 0, "initially the error count should be zero")
	testService.AddFailure()
	assert.Equal(suite.T(), testService.Fails, 1, "initially the error count should be zero")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTheTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
