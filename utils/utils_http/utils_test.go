package httputils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonPrettyPrintNoError(t *testing.T) {
	ugly := `{"one":"one", "two":"two", "three":{"four":"four"}}`
	beauty := JSONPrettyPrint(ugly)
	assert.EqualValues(t, `{
	"one": "one",
	"two": "two",
	"three": {
		"four": "four"
	}
}`, beauty)
}

func TestJsonPrettyPrintBadJson(t *testing.T) {
	ugly := `{"one":"one", "two":"two", "three":{`
	beauty := JSONPrettyPrint(ugly)
	assert.EqualValues(t, `{"one":"one", "two":"two", "three":{`, beauty)
}
