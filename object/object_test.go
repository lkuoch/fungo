package object

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ObjectTestSuite struct {
	suite.Suite
}

func TestObjectTestSuite(t *testing.T) {
	suite.Run(t, &ObjectTestSuite{})
}

func (t *ObjectTestSuite) TestStringHashKey() {
	test1a := &String{Value: "Hello World"}
	test1b := &String{Value: "Hello World"}

	test2a := &String{Value: "My name is law"}
	test2b := &String{Value: "My name is law"}

	t.Equal(test1a.HashKey(), test1b.HashKey())
	t.Equal(test2a.HashKey(), test2b.HashKey())
	t.NotEqual(test1a.HashKey(), test2a.HashKey())
	t.NotEqual(test1b.HashKey(), test2b.HashKey())
}
