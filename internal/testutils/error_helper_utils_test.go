package testutils

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestContainsError(t *testing.T) {

	t.Run("Detect Wrapped error", func(t *testing.T) {
		assert.False(t,
			ContainsError(errors.Wrap(errors.New("Contained"), "Wrapping"), errors.New("Other Error")),
		)
	})

	t.Run("Detect Wrapped error", func(t *testing.T) {
		assert.True(t,
			ContainsError(errors.Wrap(errors.New("Contained"), "Wrapping"), errors.New("Contained")),
		)
	})

	t.Run("Detect Wrapping error", func(t *testing.T) {
		assert.True(t,
			ContainsError(errors.Wrap(errors.New("Contained"), "Wrapping"), errors.New("Wrapping")),
		)
	})

	t.Run("Detect regular error", func(t *testing.T) {
		assert.True(t,
			ContainsError(errors.New("Error"), errors.New("Error")),
		)
	})
}
