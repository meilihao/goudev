package goudev

// #include <libudev.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	ErrDeviceNotFoundByName = errors.New("DeviceNotFoundByName")
	ErrNoParentDevice       = errors.New("NoParentDevice")
)

type Device struct {
	udevDevice *C.struct_udev_device
}

func (d *Device) Free() {
	if d.udevDevice != nil {
		C.udev_device_unref(d.udevDevice)
	}
}

func NewDevice() *Device {
	return &Device{}
}

func (d *Device) FromName(ctx *Context, subsystem, sysName string) error {
	cSubsystem := C.CString(subsystem)
	defer C.free(unsafe.Pointer(cSubsystem))

	cSysName := C.CString(sysName)
	defer C.free(unsafe.Pointer(cSysName))

	d.udevDevice = C.udev_device_new_from_subsystem_sysname(ctx.udev, cSubsystem, cSysName)
	if d.udevDevice == nil {
		return ErrDeviceNotFoundByName
	}

	return nil
}

func (d *Device) SysPath() string {
	return C.GoString(C.udev_device_get_syspath(d.udevDevice))
}

func (d *Device) String() string {
	return fmt.Sprintf(`Device("%s")`, d.SysPath())
}

func (d *Device) Get(property string) string {
	cProperty := C.CString(property)
	defer C.free(unsafe.Pointer(cProperty))

	return C.GoString(C.udev_device_get_property_value(d.udevDevice, cProperty))
}

func (d *Device) GetAttribute(attribute string) string {
	cAttribute := C.CString(attribute)
	defer C.free(unsafe.Pointer(cAttribute))

	return C.GoString(C.udev_device_get_sysattr_value(d.udevDevice, cAttribute))
}

type ListEntry struct {
	Name  string
	Value string
}

type ListEntryArray []ListEntry

func NewListEntryArray(ptr *C.struct_udev_list_entry) ListEntryArray {
	le := ListEntryArray{}
	for ptr != nil {
		le = append(le, ListEntry{
			Name:  C.GoString(C.udev_list_entry_get_name(ptr)),
			Value: C.GoString(C.udev_list_entry_get_value(ptr)), // udev_device_get_properties_list_entry have, udev_device_get_sysattr_list_entry no value.
		})

		ptr = C.udev_list_entry_get_next(ptr)
	}
	return le
}

type ListEntryMap map[string]string

func (d *Device) Properties() ListEntryMap {
	entries := NewListEntryArray(C.udev_device_get_properties_list_entry(d.udevDevice))
	properties := ListEntryMap{}
	for _, entry := range entries {
		properties[entry.Name] = entry.Value // d.Get(entry.Name)
	}
	return properties
}

func (d *Device) Attributes() ListEntryMap {
	entries := NewListEntryArray(C.udev_device_get_sysattr_list_entry(d.udevDevice))
	attributes := ListEntryMap{}
	for _, entry := range entries {
		attributes[entry.Name] = d.GetAttribute(entry.Name)
	}
	return attributes
}

func (d *Device) Parent() (*Device, error) {
	p := C.udev_device_get_parent(d.udevDevice)
	if p == nil {
		return nil, ErrNoParentDevice
	}

	// the parent device is not referenced, thus forcibly acquire a reference
	return &Device{
		udevDevice: C.udev_device_ref(p),
	}, nil
}

type FilterFn func(td *Device) bool

func WithFilterPciParentChildren(p *Device) func(*Device) bool {
	return func(td *Device) bool {
		if p.SysPath() == td.SysPath() {
			return false
		}

		if td.Get("PCI_SLOT_NAME") == "" || td.Get("PCI_ID") == "" {
			return false
		}

		return true
	}
}

func (d *Device) Children() ([]*Device, error) {
	e := &Enumerate{
		udevEnumerate: C.udev_enumerate_new(C.udev_device_get_udev(d.udevDevice)),
	}
	defer e.Free()

	err := e.MatchParent(d)
	if err != nil {
		return nil, err
	}

	return e.Devices(WithFilterPciParentChildren(d))
}
