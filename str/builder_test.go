package str

import (
	"github.com/nicholaskh/assert"
	"testing"
)

func TestStringBuilder(t *testing.T) {
	sb := NewStringBuilder()
	sb.WriteString("hello")
	sb.WriteString("world")
	assert.Equal(t, "helloworld", sb.String())
}
