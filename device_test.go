package goudev

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestDeviceFromName(t *testing.T) {
	ctx := NewContext()
	defer ctx.Free()

	d := NewDevice()
	defer d.Free()

	err := d.FromName(ctx, "pci", "0000:65:00.0")
	assert.Nil(t, err)

	spew.Dump(d.String())

	vendor := d.Get("ID_VENDOR_FROM_DATABASE")
	assert.NotEmpty(t, vendor)
	spew.Dump(vendor)

	vendor2 := d.Get("ID_VENDOR_FROM_DATABASE2")
	assert.Empty(t, vendor2)
	spew.Dump(vendor2)

	props := d.Properties()
	spew.Dump(props)

	attrs := d.Attributes()
	spew.Dump(attrs)

	p, err := d.Parent()
	assert.Nil(t, err)
	defer p.Free()

	cs, err := p.Children()
	assert.Nil(t, err)
	for _, cd := range cs {
		spew.Dump(cd)
		cd.Free()
	}
}
