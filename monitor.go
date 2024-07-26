package goudev

// #include <libudev.h>
// #include <stdlib.h>
import "C"
import (
	"context"
	"errors"
	"unsafe"

	"golang.org/x/sys/unix"
)

type Monitor struct {
	udevMoniter *C.struct_udev_monitor
}

func (m *Monitor) Free() {
	if m.udevMoniter != nil {
		C.udev_monitor_unref(m.udevMoniter)
	}
}

func (m *Monitor) SetReceiveBufferSize(size int) error {
	if C.udev_monitor_set_receive_buffer_size(m.udevMoniter, (C.int)(size)) != 0 {
		return errors.New("udev: udev_monitor_set_receive_buffer_size failed")
	}

	return nil
}

func (m *Monitor) FilterBy(subsystem, deviceType string) error {
	cSubsystem := C.CString(subsystem)
	defer C.free(unsafe.Pointer(cSubsystem))

	var cDeviceType *C.char
	if deviceType != "" {
		cDeviceType = C.CString(deviceType)
		defer C.free(unsafe.Pointer(cDeviceType))
	}

	if C.udev_monitor_filter_add_match_subsystem_devtype(m.udevMoniter, cSubsystem, cDeviceType) != 0 {
		return errors.New("udev: udev_monitor_filter_add_match_subsystem_devtype failed")
	}

	if C.udev_monitor_filter_update(m.udevMoniter) != 0 {
		return errors.New("udev: udev_monitor_filter_update failed")
	}

	return nil
}

func (m *Monitor) FilterByTag(tag string) error {
	cTag := C.CString(tag)
	defer C.free(unsafe.Pointer(cTag))

	if C.udev_monitor_filter_add_match_tag(m.udevMoniter, cTag) != 0 {
		return errors.New("udev: udev_monitor_filter_add_match_tag failed")
	}

	if C.udev_monitor_filter_update(m.udevMoniter) != 0 {
		return errors.New("udev: udev_monitor_filter_update failed")
	}

	return nil
}

func (m *Monitor) RemoveFilter() error {
	if C.udev_monitor_filter_remove(m.udevMoniter) != 0 {
		return errors.New("udev: udev_monitor_filter_remove failed")
	}

	if C.udev_monitor_filter_update(m.udevMoniter) != 0 {
		return errors.New("udev: udev_monitor_filter_update failed")
	}

	return nil
}

// todo
func (m *Monitor) DeviceChan(ctx context.Context) (<-chan *Device, error) {
	if C.udev_monitor_enable_receiving(m.udevMoniter) != 0 {
		return nil, errors.New("udev: udev_monitor_enable_receiving failed")
	}

	// Force monitor FD into non-blocking mode
	fd := C.udev_monitor_get_fd(m.udevMoniter)
	if e := unix.SetNonblock(int(fd), true); e != nil {
		return nil, errors.New("udev: unix.SetNonblock failed")
	}

	return nil, nil
}
