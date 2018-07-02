// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package d2xx is a cgo wrapper for the FTDI d2xx drive.
//
// The supported devices (ft232h/ft232r) implement support for various
// protocols like the GPIO, I²C, SPI, UART, JTAG.
//
// More details
//
// See https://periph.io/device/ftdi/ for more details, and how to configure
// the host to be able to use this driver.
//
// Datasheets
//
// http://www.ftdichip.com/Support/Documents/DataSheets/ICs/DS_FT232R.pdf
//
// http://www.ftdichip.com/Support/Documents/DataSheets/ICs/DS_FT232H.pdf
package d2xx
