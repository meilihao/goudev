package goudev

import (
	"path/filepath"
	"strings"
)

var (
	Devices = new(devices)
)

type devices struct{}

func (ds *devices) FromName(ctx *Context, subsystem, sysName string) (*Device, error) {
	dev := NewDevice()

	if err := dev.FromName(ctx, subsystem, sysName); err != nil {
		return nil, err
	}

	return dev, nil
}

func (ds *devices) FromPath(ctx *Context, path string) (*Device, error) {
	if !strings.HasPrefix(path, "/sys") {
		path = filepath.Join("/sys", path)
	}

	return ds.FromSysPath(ctx, path)
}

func (ds *devices) FromSysPath(ctx *Context, path string) (*Device, error) {
	dev := NewDevice()

	if err := dev.FromSysPath(ctx, path); err != nil {
		return nil, err
	}

	return dev, nil
}
