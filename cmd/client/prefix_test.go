package client

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestPrefixParsing(t *testing.T) {
	text := "!play"
	text2 := "!pause"

	assert.True(t, strings.HasPrefix(text, "!play"))
	assert.True(t, strings.HasPrefix(text2, "!pause"))
}