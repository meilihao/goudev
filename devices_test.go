package goudev

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFromName(t *testing.T) {
	ctx := NewContext()
	defer ctx.Free()

	d, err := Devices.FromName(ctx, "pci", "0000:65:00.0")
	assert.Nil(t, err)
	defer d.Free()

	spew.Dump(d.String())
}
