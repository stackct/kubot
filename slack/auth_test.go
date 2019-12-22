package slack

import (
	"kubot/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasAccess(t *testing.T) {
	config.Conf = config.Config{
		Environments: []config.Environment{
			config.Environment{Name: "e1", Users: []string{"u1"}},
			config.Environment{Name: "e2", Users: []string{"u2", "u3"}},
		},
	}

	testCases := []struct {
		environment string
		user        string
		expected    bool
	}{
		{"e1", "u1", true},
		{"e1", "u1", true},
		{"e2", "u2", true},
		{"e2", "u3", true},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, HasAccess(tc.user, tc.environment))
	}
}
