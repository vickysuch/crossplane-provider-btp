package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVal(t *testing.T) {

	// nil pointer
	var ptrString *string
	assert.Equal(t, "", Val(ptrString))

	// value pointer
	str := "Foo"
	ptrString = &str
	assert.Equal(t, "Foo", Val(ptrString))

	// pointer to empty value
	emptyStr := ""
	ptrString = &emptyStr
	assert.Equal(t, "", Val(ptrString))

	// same tests for bool to ensure its generic
	var ptrBool *bool
	assert.Equal(t, false, Val(ptrBool))
	b := true
	ptrBool = &b
	assert.Equal(t, true, Val(ptrBool))

	emptyB := false
	ptrBool = &emptyB
	assert.Equal(t, false, Val(ptrBool))

}
