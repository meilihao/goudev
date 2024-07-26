//go:build linux && cgo

package goudev

// #cgo pkg-config: udev
// #include <libudev.h>
import "C"

type Context struct {
	udev *C.struct_udev
}

func NewContext() *Context {
	return &Context{
		udev: C.udev_new(),
	}
}

func (c *Context) Free() {
	if c.udev != nil {
		C.udev_unref(c.udev)
	}
}

func (c *Context) NewEnumerate() *Enumerate {
	return &Enumerate{
		udevEnumerate: C.udev_enumerate_new(c.udev),
	}
}
