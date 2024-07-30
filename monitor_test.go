package goudev

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/meilihao/golib/v2/cmd"
	"github.com/stretchr/testify/assert"
)

func TestMonitor(t *testing.T) {
	c := NewContext()
	defer c.Free()

	m := c.NewMonitor()
	defer m.Free()

	m.FilterBy("block")
	m.FilterBy("dlm")
	m.FilterBy("net")

	ctx, cancel := context.WithCancel(context.Background())
	ch, e := m.DeviceChan(ctx, 0)
	assert.Nil(t, e)

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		cmd.CmdCombinedBash(nil, "qemu-nbd -c /dev/nbd0 disk.img")
		cmd.CmdCombinedBash(nil, "qemu-nbd -d /dev/nbd0")
		wg.Done()
	}()
	go func() {
		fmt.Println("Started listening on channel")
		for d := range ch {
			fmt.Println(d.SysPath(), d.Action())
			d.Free()
		}
		fmt.Println("Channel closed")
		wg.Done()
	}()
	go func() {
		fmt.Println("Starting timer to Removing filter done")
		<-time.After(2 * time.Second)
		fmt.Println("Removing filter")
		m.RemoveFilter()
		wg.Done()
	}()
	go func() {
		fmt.Println("Starting timer to signal done")
		<-time.After(4 * time.Second)
		fmt.Println("Signalling done")
		cancel()
		wg.Done()
	}()
	wg.Wait()
}
