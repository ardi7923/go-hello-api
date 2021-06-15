package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	// assert equality
	assert.Equal(t, 123, 123, "they should be equal")

}
