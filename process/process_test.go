package process

import (
	"testing"
)

func TestStart(t *testing.T) {
	Start("echo", []string{}, map[string]string{"foo3": "bar3"})
}
