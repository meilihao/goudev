package goudev

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestDeviceFromNamePci(t *testing.T) {
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

func TestDeviceFromNameBlock(t *testing.T) {
	ctx := NewContext()
	defer ctx.Free()

	d := NewDevice()
	defer d.Free()

	err := d.FromName(ctx, "block", "nvme0n1")
	assert.Nil(t, err)

	spew.Dump(d.String())
}

func TestDeviceFromPath(t *testing.T) {
	ctx := NewContext()
	defer ctx.Free()

	d := NewDevice()
	defer d.Free()

	err := d.FromPath(ctx, "class/block/nvme0n1")
	assert.Nil(t, err)

	spew.Dump(d.String())
}

func TestDeviceFromSysPath(t *testing.T) {
	ctx := NewContext()
	defer ctx.Free()

	d := NewDevice()
	defer d.Free()

	err := d.FromPath(ctx, "/sys/class/block/nvme0n1")
	assert.Nil(t, err)

	spew.Dump(d.String())
	spew.Dump(d.SysName())
}

func TestDeviceAttributes(t *testing.T) {
	ctx := NewContext()
	defer ctx.Free()

	d := NewDevice()
	defer d.Free()

	err := d.FromPath(ctx, "/sys/class/block/nvme0n1")
	assert.Nil(t, err)

	spew.Dump(d.String())
	attrs := d.Attributes()
	spew.Dump(attrs)
	spew.Dump(attrs["size"])
	spew.Dump(attrs["queue/logical_block_size"])
	spew.Dump(attrs["queue/rotational"])
	spew.Dump(attrs["protection_type"])
}

func TestDeviceProperties(t *testing.T) {
	ctx := NewContext()
	defer ctx.Free()

	d := NewDevice()
	defer d.Free()

	err := d.FromPath(ctx, "/sys/class/block/nvme0n1")
	assert.Nil(t, err)

	spew.Dump(d.String())
	props := d.Properties()
	spew.Dump(props)
	spew.Dump(props["ID_MODEL"])
	spew.Dump(props["DRIVER"])
	spew.Dump(props["SUBSYSTEM"])
	spew.Dump(props["DEVPATH"])
	spew.Dump(props["ID_WWN"])
	spew.Dump(props["ID_BUS"])
}
