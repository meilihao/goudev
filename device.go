package goudev

// #include <libudev.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"unsafe"
)

var (
	ErrDeviceNotFoundByNameTmpl = "No device %s in %s"
	ErrDeviceNotFoundByPathTmpl = "No device at %s"
	ErrNoParentDevice           = errors.New("NoParentDevice")
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
		return fmt.Errorf(ErrDeviceNotFoundByNameTmpl, sysName, subsystem)
	}

	return nil
}

func (d *Device) FromPath(ctx *Context, path string) error {
	if !strings.HasPrefix(path, "/sys") {
		path = filepath.Join("/sys", path)
	}

	return d.FromSysPath(ctx, path)
}

func (d *Device) FromSysPath(ctx *Context, path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	d.udevDevice = C.udev_device_new_from_syspath(ctx.udev, cPath)
	if d.udevDevice == nil {
		return fmt.Errorf(ErrDeviceNotFoundByPathTmpl, path)
	}

	return nil
}

func (d *Device) Action() string {
	return C.GoString(C.udev_device_get_action(d.udevDevice))
}

func (d *Device) DeviceNode() string {
	return C.GoString(C.udev_device_get_devnode(d.udevDevice))
}

func (d *Device) DeviceNumber() *Devnum {
	return &Devnum{C.udev_device_get_devnum(d.udevDevice)}
}

// devpath
func (d *Device) DevicePath() string {
	return C.GoString(C.udev_device_get_devpath(d.udevDevice))
}

func (d *Device) DeviceType() string {
	return C.GoString(C.udev_device_get_devtype(d.udevDevice))
}

func (d *Device) Driver() string {
	return C.GoString(C.udev_device_get_driver(d.udevDevice))
}

func (d *Device) HasTag(tag string) bool {
	cTag := C.CString(tag)
	defer C.free(unsafe.Pointer(cTag))

	return C.udev_device_has_tag(d.udevDevice, cTag) != 0
}

func (d *Device) IsInitialized() bool {
	return C.udev_device_get_is_initialized(d.udevDevice) != 0
}

func (d *Device) TimeSinceInitialized() uint64 {
	return uint64(C.udev_device_get_usec_since_initialized(d.udevDevice)) // ms
}

func (d *Device) SequenceNumber() uint64 {
	return uint64(C.udev_device_get_seqnum(d.udevDevice))
}

func (d *Device) SysPath() string {
	return C.GoString(C.udev_device_get_syspath(d.udevDevice))
}

func (d *Device) SysName() string {
	return C.GoString(C.udev_device_get_sysname(d.udevDevice))
}

func (d *Device) SysNumber() string {
	return C.GoString(C.udev_device_get_sysnum(d.udevDevice))
}

func (d *Device) String() string {
	return fmt.Sprintf(`Device("%s")`, d.SysPath())
}

func (d *Device) Subsystem() string {
	return C.GoString(C.udev_device_get_subsystem(d.udevDevice))
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

func (d *Device) SetAttribute(sysattr, value string) error {
	cSysattr := C.CString(sysattr)
	defer C.free(unsafe.Pointer(cSysattr))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	if C.udev_device_set_sysattr_value(d.udevDevice, cSysattr, cValue) != 0 {
		return errors.New("udev: udev_device_get_sysattr_value failed")
	}

	return nil
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

func (d *Device) Attributes() ListEntryMap {
	entries := NewListEntryArray(C.udev_device_get_sysattr_list_entry(d.udevDevice))
	attributes := ListEntryMap{}
	for _, entry := range entries {
		attributes[entry.Name] = d.GetAttribute(entry.Name)
	}
	return attributes
}

// devlinks
func (d *Device) DeviceLinks() ListEntryMap {
	entries := NewListEntryArray(C.udev_device_get_devlinks_list_entry(d.udevDevice))
	devlinks := ListEntryMap{}
	for _, entry := range entries {
		devlinks[entry.Name] = entry.Value // d.Get(entry.Name)
	}
	return devlinks
}

func (d *Device) Properties() ListEntryMap {
	entries := NewListEntryArray(C.udev_device_get_properties_list_entry(d.udevDevice))
	properties := ListEntryMap{}
	for _, entry := range entries {
		properties[entry.Name] = entry.Value // d.Get(entry.Name)
	}
	return properties
}

func (d *Device) Tags() ListEntryMap {
	entries := NewListEntryArray(C.udev_device_get_tags_list_entry(d.udevDevice))
	tags := ListEntryMap{}
	for _, entry := range entries {
		tags[entry.Name] = entry.Value // d.Get(entry.Name)
	}
	return tags
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

func (d *Device) FindParent(subsystem string, deviceType ...string) (*Device, error) {
	cSubsystem := C.CString(subsystem)
	defer C.free(unsafe.Pointer(cSubsystem))

	var cDeviceType *C.char
	if len(deviceType) > 0 {
		cDeviceType = C.CString(deviceType[0])
		defer C.free(unsafe.Pointer(cDeviceType))
	} else {
		cDeviceType = nil
	}

	p := C.udev_device_get_parent_with_subsystem_devtype(d.udevDevice, cSubsystem, cDeviceType)
	if p == nil {
		return nil, ErrNoParentDevice
	}

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

func WithFilterBlockDevtype(devtype string) func(*Device) bool {
	return func(td *Device) bool {
		return td.Get("DEVTYPE") == devtype
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
