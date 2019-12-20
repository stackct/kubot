package config

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParseFile(t *testing.T) {
	config, _ := ParseFile("./resources/authorization.yml")

	expected := Config{
		Environments: []Environment{
			Environment{Name: "e1", Groups: []string{"all"}, Users: []string{}},
			Environment{Name: "e2", Groups: []string{}, Users: []string{"john.doe", "mary.foo"}},
		},
	}

	assert.True(t, reflect.DeepEqual(expected, config))
}
