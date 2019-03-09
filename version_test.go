package bitcask

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullVersion(t *testing.T) {
	assert := assert.New(t)

	expected := fmt.Sprintf("%s@%s", Version, Commit)
	assert.Equal(expected, FullVersion())
}
