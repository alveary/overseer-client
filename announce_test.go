package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AnnounceTestSuite struct {
	suite.Suite
}

func (suite *AnnounceTestSuite) SetupTest() {
}

func (suite *AnnounceTestSuite) TestAnnouncement() {
	// assert.True(suite.T(), true, "True should be true")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAnnounceTestSuite(t *testing.T) {
	suite.Run(t, new(AnnounceTestSuite))
}
