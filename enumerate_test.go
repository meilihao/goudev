package goudev

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestEnumerate(t *testing.T) {
	ctx := NewContext()
	defer ctx.Free()

	d := NewDevice()
	defer d.Free()

	err := d.FromName(ctx, "pci", "0000:65:00.0")
	assert.Nil(t, err)

	spew.Dump(d.String())

	p, err := d.Parent()
	assert.Nil(t, err)
	defer p.Free()
	spew.Dump(p.String())

	e := ctx.NewEnumerate()
	defer e.Free()

	err = e.MatchParent(p)
	assert.Nil(t, err)

	ds, err := e.Devices(nil)
	assert.Nil(t, err)
	spew.Dump(ds)
}

func TestEnumerateMatchSubsystem(t *testing.T) {
	ctx := NewContext()
	defer ctx.Free()

	e := ctx.NewEnumerate()
	defer e.Free()

	err := e.MatchSubsystem("block")
	assert.Nil(t, err)

	ds, err := e.Devices(nil)
	assert.Nil(t, err)
	spew.Dump(ds)

	dsDisk, err := e.Devices(WithFilterBlockDevtype("disk"))
	assert.Nil(t, err)
	spew.Dump(dsDisk)
}
