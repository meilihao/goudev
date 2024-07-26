package goudev

// #cgo LDFLAGS: -ludev
// #include <libudev.h>
import "C"
import (
	"errors"
)

type Enumerate struct {
	udevEnumerate *C.struct_udev_enumerate
}

func (e *Enumerate) Free() {
	if e.udevEnumerate != nil {
		C.udev_enumerate_unref(e.udevEnumerate)
	}
}

// https://github.com/systemd/systemd/blob/main/src/libudev/libudev-enumerate.c#L329
func (e *Enumerate) MatchParent(parent *Device) error {
	if C.udev_enumerate_add_match_parent(e.udevEnumerate, parent.udevDevice) != 0 {
		return errors.New("udev: udev_enumerate_add_match_parent failed")
	}

	return nil
}

// https://github.com/systemd/systemd/blob/main/src/libudev/libudev-enumerate.c#L440
func (e *Enumerate) Devices(filter FilterFn) (m []*Device, err error) {
	if C.udev_enumerate_scan_devices(e.udevEnumerate) != 0 {
		err = errors.New("udev: udev_enumerate_scan_devices failed")
	} else {
		m = make([]*Device, 0)
		for l := C.udev_enumerate_get_list_entry(e.udevEnumerate); l != nil; l = C.udev_list_entry_get_next(l) {
			s := C.udev_list_entry_get_name(l)

			d := &Device{
				udevDevice: C.udev_device_new_from_syspath(C.udev_enumerate_get_udev(e.udevEnumerate), s),
			}

			if filter != nil {
				if !filter(d) {
					d.Free()
					continue
				}
			}

			m = append(m, d)
		}
	}
	return
}
