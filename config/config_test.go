package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	file, _ := os.Open("./resources/authorization.yml")
	input, _ := ioutil.ReadAll(file)
	bytes := []byte(input)

	config, _ := Parse(bytes)

	expected := Config{
		Environments: []Environment{
			Environment{Name: "e1", Groups: []string{"all"}, Users: []string{}},
			Environment{Name: "e2", Groups: []string{}, Users: []string{"john.doe", "mary.foo"}},
		},
	}

	assert.True(t, reflect.DeepEqual(expected, config))
}
