// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/edgexfoundry/device-sdk-go"
	"github.com/edgexfoundry/device-rfrain/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-rfrain"
)

func main() {
	d := driver.SimpleDriver{}
	startup.Bootstrap(serviceName, device.Version, &d)
}
