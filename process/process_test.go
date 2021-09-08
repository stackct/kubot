package process

import (
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/memory"
	"github.com/stretchr/testify/assert"
)

func TestRunProcessWithoutInterpolation(t *testing.T) {
	h := memory.New()
	log.SetHandler(h)

	Start("echo", []string{"foo"}, map[string]string{}, map[string]string{})

	assert.Equal(t, 2, len(h.Entries))
	e := h.Entries[0]
	assert.Equal(t, e.Message, "executing command")
	assert.Equal(t, e.Level, log.InfoLevel)
	assert.Equal(t, log.Fields{"name": "echo", "args": []string{"foo"}}, e.Fields)

	e = h.Entries[1]
	assert.Equal(t, e.Message, "command completed")
	assert.Equal(t, e.Level, log.InfoLevel)
	assert.Equal(t, log.Fields{"args": []string{"foo"}, "name": "echo", "output": "foo\n"}, e.Fields)
}

func TestRunProcessWithInterpolation(t *testing.T) {
	h := memory.New()
	log.SetHandler(h)

	Start("echo", []string{"${foo}"}, map[string]string{"foo": "bar"}, map[string]string{})

	assert.Equal(t, 2, len(h.Entries))
	e := h.Entries[0]
	assert.Equal(t, e.Message, "executing command")
	assert.Equal(t, e.Level, log.InfoLevel)
	assert.Equal(t, log.Fields{"name": "echo", "args": []string{"bar"}}, e.Fields)

	e = h.Entries[1]
	assert.Equal(t, e.Message, "command completed")
	assert.Equal(t, e.Level, log.InfoLevel)
	assert.Equal(t, log.Fields{"args": []string{"bar"}, "name": "echo", "output": "bar\n"}, e.Fields)
}

func TestRunProcessWithFailure(t *testing.T) {
	h := memory.New()
	log.SetHandler(h)

	Start("commandnotfound", []string{}, map[string]string{}, map[string]string{})

	assert.Equal(t, 2, len(h.Entries))
	e := h.Entries[0]
	assert.Equal(t, e.Message, "executing command")
	assert.Equal(t, e.Level, log.InfoLevel)
	assert.Equal(t, log.Fields{"name": "commandnotfound", "args": []string{""}}, e.Fields)

	e = h.Entries[1]
	assert.Equal(t, e.Message, "command failed")
	assert.Equal(t, e.Level, log.ErrorLevel)
}

func TestRunProcessCapturesStdOutAndStdErr(t *testing.T) {
	h := memory.New()
	log.SetHandler(h)

	Start("sh", []string{"-c", "echo stdout; echo 1>&2 stderr"}, map[string]string{}, map[string]string{})

	assert.Equal(t, 2, len(h.Entries))
	e := h.Entries[0]
	assert.Equal(t, e.Message, "executing command")
	assert.Equal(t, e.Level, log.InfoLevel)

	e = h.Entries[1]
	assert.Equal(t, e.Message, "command completed")
	assert.Equal(t, e.Fields["output"], "stdout\nstderr\n")
}
