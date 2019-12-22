package command

import (
	"kubot/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdoutWriter(t *testing.T) {
	n, err := StdoutWriter{}.Write([]byte("foo"))
	assert.Nil(t, err)
	assert.Equal(t, 3, n)
}

func TestStderrWriter(t *testing.T) {
	n, err := StderrWriter{}.Write([]byte("foo"))
	assert.Nil(t, err)
	assert.Equal(t, 3, n)
}

func TestExecuteFailureWhenCommandDoesNotExist(t *testing.T) {
	err := Execute("foo", map[string]string{})
	assert.Equal(t, "command not found: foo", err.Error())
}

func TestExecuteFailureWhenCommandFails(t *testing.T) {
	config.Conf = config.NewMockConfig()
	err := Execute("fail", map[string]string{})
	assert.Equal(t, "exit status 2", err.Error())
}

func TestExecuteSuccess(t *testing.T) {
	config.Conf = config.NewMockConfig()
	err := Execute("help", map[string]string{})
	assert.Nil(t, err)
}

func TestExecuteInterpolatesFromCommandConfig(t *testing.T) {
	config.Conf = config.NewMockConfig()
	mockWriter := &config.MockWriter{}
	CommandStdoutWriter = mockWriter
	err := Execute("echo", map[string]string{"foo3": "bar3"})
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo bar1 bar2 bar3\n"}, mockWriter.Messages)
}
