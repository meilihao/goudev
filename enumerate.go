package goudev

// #include <libudev.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"
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

// https://github.com/systemd/systemd/blob/main/src/libudev/libudev-enumerate.c#L278
func (e *Enumerate) MatchProperty(prop, value string) error {
	cProp := C.CString(prop)
	defer C.free(unsafe.Pointer(cProp))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	if C.udev_enumerate_add_match_property(e.udevEnumerate, cProp, cValue) != 0 {
		return errors.New("udev: udev_enumerate_add_match_property failed")
	}

	return nil
}

func (e *Enumerate) MatchSysattr(sysattr, value string) error {
	cSysattr := C.CString(sysattr)
	defer C.free(unsafe.Pointer(cSysattr))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	if C.udev_enumerate_add_match_sysattr(e.udevEnumerate, cSysattr, cValue) != 0 {
		return errors.New("udev: udev_enumerate_add_match_sysattr failed")
	}

	return nil
}

// https://github.com/systemd/systemd/blob/main/src/libudev/libudev-enumerate.c
func (e *Enumerate) MatchSubsystem(subsystem string) error {
	cSubsystem := C.CString(subsystem)
	defer C.free(unsafe.Pointer(cSubsystem))

	if C.udev_enumerate_add_match_subsystem(e.udevEnumerate, cSubsystem) != 0 {
		return errors.New("udev: udev_enumerate_add_match_subsystem failed")
	}

	return nil
}

// https://github.com/systemd/systemd/blob/main/src/libudev/libudev-enumerate.c#L385
func (e *Enumerate) MatchSysname(sysname string) error {
	cSysname := C.CString(sysname)
	defer C.free(unsafe.Pointer(cSysname))

	if C.udev_enumerate_add_match_sysname(e.udevEnumerate, cSysname) != 0 {
		return errors.New("udev: udev_enumerate_add_match_sysname failed")
	}

	return nil
}

// https://github.com/systemd/systemd/blob/main/src/libudev/libudev-enumerate.c#L303
func (e *Enumerate) MatchTag(tag string) error {
	cTag := C.CString(tag)
	defer C.free(unsafe.Pointer(cTag))

	if C.udev_enumerate_add_match_tag(e.udevEnumerate, cTag) != 0 {
		return errors.New("udev: udev_enumerate_add_match_tag failed")
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
