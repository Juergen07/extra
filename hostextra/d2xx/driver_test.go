// Copyright 2016 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package d2xx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"periph.io/x/extra/hostextra/d2xx/ftdi"
)

func TestDriver(t *testing.T) {
	defer reset(t)
	drv.numDevices = func() (int, error) {
		return 1, nil
	}
	drv.d2xxOpen = func(i int) (d2xxHandle, int) {
		if i != 0 {
			t.Fatalf("unexpected index %d", i)
		}
		d := &d2xxFakeHandle{
			d:    ftdi.FT232R,
			vid:  0x0403,
			pid:  0x6014,
			data: [][]byte{{}, {0}},
		}
		return d, 0
	}
	if b, err := drv.Init(); !b || err != nil {
		t.Fatalf("Init() = %t, %v", b, err)
	}
}

//

type d2xxFakeHandle struct {
	d       ftdi.DevType
	vid     uint16
	pid     uint16
	data    [][]byte
	ua      []byte
	e       ftdi.EEPROM
	touched bool
}

func (d *d2xxFakeHandle) d2xxClose() int {
	return 0
}
func (d *d2xxFakeHandle) d2xxResetDevice() int {
	d.touched = true
	return 0
}
func (d *d2xxFakeHandle) d2xxGetDeviceInfo() (ftdi.DevType, uint16, uint16, int) {
	return d.d, d.vid, d.pid, 0
}
func (d *d2xxFakeHandle) d2xxEEPROMRead(dev ftdi.DevType, e *ftdi.EEPROM) int {
	d.touched = true
	*e = d.e
	return 0
}
func (d *d2xxFakeHandle) d2xxEEPROMProgram(e *ftdi.EEPROM) int {
	d.touched = true
	d.e = *e
	return 0
}
func (d *d2xxFakeHandle) d2xxEraseEE() int {
	d.touched = true
	return 0
}
func (d *d2xxFakeHandle) d2xxWriteEE(offset uint8, value uint16) int {
	d.touched = true
	return 1
}
func (d *d2xxFakeHandle) d2xxEEUASize() (int, int) {
	d.touched = true
	return len(d.ua), 0
}
func (d *d2xxFakeHandle) d2xxEEUARead(ua []byte) int {
	d.touched = true
	copy(ua, d.ua)
	return 0
}
func (d *d2xxFakeHandle) d2xxEEUAWrite(ua []byte) int {
	d.touched = true
	d.ua = make([]byte, len(ua))
	copy(d.ua, ua)
	return 0
}
func (d *d2xxFakeHandle) d2xxSetChars(eventChar byte, eventEn bool, errorChar byte, errorEn bool) int {
	d.touched = true
	return 0
}
func (d *d2xxFakeHandle) d2xxSetUSBParameters(in, out int) int {
	d.touched = true
	return 0
}
func (d *d2xxFakeHandle) d2xxSetFlowControl() int {
	d.touched = true
	return 0
}
func (d *d2xxFakeHandle) d2xxSetTimeouts(readMS, writeMS int) int {
	d.touched = true
	return 0
}
func (d *d2xxFakeHandle) d2xxSetLatencyTimer(delayMS uint8) int {
	d.touched = true
	return 0
}
func (d *d2xxFakeHandle) d2xxSetBaudRate(hz uint32) int {
	d.touched = true
	return 0
}
func (d *d2xxFakeHandle) d2xxGetQueueStatus() (uint32, int) {
	d.touched = true
	if len(d.data) == 0 {
		return 0, 0
	}
	// This is to work around flushPending().
	l := len(d.data[0])
	if l == 0 {
		d.data = d.data[1:]
	}
	return uint32(l), 0
}
func (d *d2xxFakeHandle) d2xxRead(b []byte) (int, int) {
	d.touched = true
	if len(d.data) == 0 {
		return 0, 0
	}
	l := len(b)
	if j := len(d.data[0]); j < l {
		l = j
	}
	if l == 0 {
		d.data = d.data[1:]
		return 0, 0
	}
	copy(b, d.data[0])
	d.data[0] = d.data[0][l:]
	if len(d.data[0]) == 0 {
		d.data = d.data[1:]
	}
	return l, 0
}
func (d *d2xxFakeHandle) d2xxWrite(b []byte) (int, int) {
	d.touched = true
	return 0, 0
}
func (d *d2xxFakeHandle) d2xxGetBitMode() (byte, int) {
	d.touched = true
	return 0, 0
}
func (d *d2xxFakeHandle) d2xxSetBitMode(mask, mode byte) int {
	d.touched = true
	return 0
}

func reset(t *testing.T) {
	drv.reset()
}

func init() {
	reset(nil)
}

var dev = []d2xxFakeHandle{
	{
		d:    ftdi.FT232H,
		vid:  0x0403,
		pid:  0x6014,
		data: [][]byte{{}, {0}},
	},
	{
		d:    ftdi.FT232R,
		vid:  0x0403,
		pid:  0x6001,
		data: [][]byte{{}, {0}},
	},
	{
		d:    ftdi.FT232R,
		vid:  0x0403,
		pid:  0x6002,
		data: [][]byte{{}, {0}},
	},
}

func TestFilter1(t *testing.T) {
	defer reset(t)
	drv.numDevices = func() (int, error) {
		return len(dev), nil
	}
	drv.d2xxOpen = func(i int) (d2xxHandle, int) {
		if i < 0 || i >= len(dev) {
			t.Fatalf("unexpected index %d", i)
		}
		dev[i].touched = false
		return &dev[i], 0
	}
	// filter for first FT232R device
	DeviceFilter = []Filter{{
		Type:      FilterDeviceType(ftdi.FT232R),
		DeviceIdx: FilterDeviceIdx(-1),
	}}
	if b, err := drv.Init(); !b || err != nil {
		t.Fatalf("Init() = %t, %v", b, err)
	}
	all := All()
	assert.Equal(t, 3, len(all))
	assert.Contains(t, all[0].String(), "no match filter")
	assert.False(t, dev[0].touched)
	assert.Equal(t, all[1].String(), "FT232R(1)")
	assert.True(t, dev[1].touched)
	assert.Contains(t, all[2].String(), "no match filter")
	assert.False(t, dev[2].touched)
}

func TestFilter2(t *testing.T) {
	defer reset(t)
	drv.numDevices = func() (int, error) {
		return len(dev), nil
	}
	drv.d2xxOpen = func(i int) (d2xxHandle, int) {
		if i < 0 || i >= len(dev) {
			t.Fatalf("unexpected index %d", i)
		}
		dev[i].vid++
		dev[i].touched = false
		return &dev[i], 0
	}
	// filter for first FT232R device
	DeviceFilter = []Filter{{
		Type:      FilterDeviceType(ftdi.FT232R),
		DeviceIdx: FilterDeviceIdx(1),
	}}
	if b, err := drv.Init(); !b || err != nil {
		t.Fatalf("Init() = %t, %v", b, err)
	}
	all := All()
	assert.Equal(t, 3, len(all))
	assert.Contains(t, all[0].String(), "no match filter")
	assert.False(t, dev[0].touched)
	assert.Contains(t, all[1].String(), "no match filter")
	assert.False(t, dev[1].touched)
	assert.Equal(t, all[2].String(), "FT232R(2)")
	assert.True(t, dev[2].touched)
}
