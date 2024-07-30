package goudev

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevnumMajorMinor(t *testing.T) {
	d := MkDev(1, 8)
	assert.Equal(t, d.Major(), 1)
	assert.Equal(t, d.Minor(), 8)
}
