package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvironment_ContainsUser(t *testing.T) {
	e := Environment{Users: []string{"u1"}}

	assert.True(t, e.ContainsUser("u1"))
	assert.False(t, e.ContainsUser("u2"))
}
