package command

import (
	"kubot/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute_Failure_When_Command_Fails(t *testing.T) {
	config.Conf = config.NewMockConfig()
	out := make(chan string)
	err := Execute("fail", "product", map[string]string{}, map[string]string{}, out)
	assert.NotNil(t, err)
}

func TestExecute_Success(t *testing.T) {
	config.Conf = config.NewMockConfig()
	out := make(chan string)
	err := Execute("help", "product", map[string]string{}, map[string]string{}, out)
	assert.Nil(t, err)
}

func TestExecute_Interpolates_From_Command_Config(t *testing.T) {
	config.Conf = config.NewMockConfig()
	out := make(chan string)
	err := Execute("echo", "product", map[string]string{"foo3": "bar3"}, map[string]string{}, out)
	assert.Nil(t, err)
}

func TestExecute_With_ChannelStdout(t *testing.T) {
	config.Conf = config.NewMockConfig()
	out := make(chan string)
	go Execute("echo", "stdout", map[string]string{}, map[string]string{}, out)
	assert.Equal(t, "```foo\n```", <-out)
}
