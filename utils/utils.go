package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func ProdAssert(t *testing.T) *assert.Assertions {
	return assert.New(t)
}
