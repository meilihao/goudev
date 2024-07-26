//go:build linux && cgo

package goudev

// #cgo pkg-config: udev
// #cgo LDFLAGS: -ludev
// #include <libudev.h>
// #include <stdlib.h>
import "C"
import "unsafe"

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

// only use "udev"
func (c *Context) NewMonitor() *Monitor {
	cSource := C.CString("udev")
	defer C.free(unsafe.Pointer(cSource))

	return &Monitor{
		udevMoniter: C.udev_monitor_new_from_netlink(c.udev, cSource),
	}
}
