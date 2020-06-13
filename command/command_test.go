package command

import (
	"kubot/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute_Failure_When_Command_Fails(t *testing.T) {
	config.Conf = config.NewMockConfig()
	err := Execute("fail", "product", map[string]string{}, map[string]string{})
	assert.Equal(t, "exit status 2", err.Error())
}

func TestExecute_Success(t *testing.T) {
	config.Conf = config.NewMockConfig()
	err := Execute("help", "product", map[string]string{}, map[string]string{})
	assert.Nil(t, err)
}

func TestExecute_Interpolates_From_Command_Config(t *testing.T) {
	config.Conf = config.NewMockConfig()
	err := Execute("echo", "product", map[string]string{"foo3": "bar3"}, map[string]string{})
	assert.Nil(t, err)
}
