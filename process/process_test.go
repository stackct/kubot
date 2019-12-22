package process

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
