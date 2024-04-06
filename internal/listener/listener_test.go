//go:build all || unit

package listener

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListen(t *testing.T) {
	t.Run("Should start listening", func(t *testing.T) {
		// given
		l := NewListener()

		// when
		listen := l.Listen()

		// then
		assert.True(t, listen)
	})
}
