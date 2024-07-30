package goudev

// #include <libudev.h>
// #include <stdlib.h>
import "C"
import (
	"context"
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	maxEpollEvents = 32
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

func (m *Monitor) FilterBy(subsystem string, deviceType ...string) error {
	cSubsystem := C.CString(subsystem)
	defer C.free(unsafe.Pointer(cSubsystem))

	var cDeviceType *C.char
	if len(deviceType) > 0 {
		cDeviceType = C.CString(deviceType[0])
		defer C.free(unsafe.Pointer(cDeviceType))
	} else {
		cDeviceType = nil
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

func (m *Monitor) receiveDevice() *Device {
	d := C.udev_monitor_receive_device(m.udevMoniter)
	if d == nil {
		return nil
	}

	return &Device{
		udevDevice: d,
	}
}

// epollTimeout, ms
func (m *Monitor) DeviceChan(ctx context.Context, epollTimeout int) (<-chan *Device, error) {
	if C.udev_monitor_enable_receiving(m.udevMoniter) != 0 {
		return nil, errors.New("udev: udev_monitor_enable_receiving failed")
	}

	// Force monitor FD into non-blocking mode
	fd := C.udev_monitor_get_fd(m.udevMoniter)
	if e := unix.SetNonblock(int(fd), true); e != nil {
		return nil, errors.New("udev: unix.SetNonblock failed")
	}

	// Create an epoll fd
	epfd, e := unix.EpollCreate1(0)
	if e != nil {
		return nil, errors.New("udev: unix.EpollCreate1 failed")
	}

	var event unix.EpollEvent
	var events [maxEpollEvents]unix.EpollEvent

	// Add the fd to the epoll fd
	event.Events = unix.EPOLLIN | unix.EPOLLET
	event.Fd = int32(fd)
	if e = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(fd), &event); e != nil {
		return nil, errors.New("udev: unix.EpollCtl failed")
	}

	// Create the channel
	ch := make(chan *Device)

	// Create goroutine to epoll the fd
	go func(fd int32) {
		// Close the epoll fd when goroutine exits
		defer unix.Close(epfd)
		// Close the channel when goroutine exits
		defer close(ch)
		// Loop forever
		for {
			// Poll the file descriptor
			nevents, e := unix.EpollWait(epfd, events[:], epollTimeout)
			// Ignore the EINTR error case since cancelation is performed with the
			// context's Done() channel
			if e != nil {
				errno, isErrno := e.(syscall.Errno)
				if isErrno && errno == syscall.EINTR {
					continue
				} else {
					return
				}
			}

			// if (e != nil && !isErrno) || (isErrno && errno != syscall.EINTR) {
			// 	return
			// }

			// Check for done signal
			select {
			case <-ctx.Done():
				return
			default:
			}
			// Process events
			for ev := 0; ev < nevents; ev++ {
				if events[ev].Fd == fd {
					if (events[ev].Events & unix.EPOLLIN) != 0 {
						for d := m.receiveDevice(); d != nil; d = m.receiveDevice() {
							ch <- d
						}
					}
				}
			}
		}
	}(int32(fd))

	return ch, nil
}
