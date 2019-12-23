package command

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMockParser(t *testing.T) {
	cmd, err := NewMockParser().Parse("foo")

	assert.Nil(t, err)
	assert.NotNil(t, cmd)
}

func TestMockParser_With_Options(t *testing.T) {
	cmd, err := NewMockParser(func(mp *MockParser) {
		mp.MockError = errors.New("oops")
		mp.MockCommand = nil
	}).Parse("foo")

	assert.Nil(t, cmd)
	assert.Equal(t, "oops", err.Error())
}

func TestMockCommand_Name(t *testing.T) {
	c := MockCommand{}

	assert.Equal(t, "Mock", c.Name())
}

func TestMockCommand_Execute(t *testing.T) {
	out := make(chan string)
	c := MockCommand{}

	go c.Execute(out, Context{})
	assert.Equal(t, "fin", <-out)
}
