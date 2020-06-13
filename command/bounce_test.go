package command

import (
	"kubot/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBounce(t *testing.T) {
	b, err := NewBounce([]string{})

	assert.Equal(t, b, &Bounce{})
	assert.Equal(t, err, nil)
}

func TestBounce_Execute(t *testing.T) {
	out := make(chan string)
	config.Conf = config.NewMockConfig()

	go Bounce{}.Execute(out, Context{Environment: config.Environment{Name: "local"}})

	assert.Equal(t, "Bouncing now", <-out)
}
