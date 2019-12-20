package config

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParseFile(t *testing.T) {
	config, _ := ParseFile("./resources/config.yml")

	expected := Config{
		Environments: []Environment{
			Environment{Name: "e1", Users: []string{}},
			Environment{Name: "e2", Users: []string{"john.doe", "mary.foo"}},
		},
	}

	assert.True(t, reflect.DeepEqual(expected, config))
}

func TestEnvironment_ContainsUser(t *testing.T) {
	e := Environment{Users: []string{"u1"}}

	assert.True(t, e.ContainsUser("u1"))
	assert.False(t, e.ContainsUser("u2"))
}
