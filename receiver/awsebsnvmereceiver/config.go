// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package awsebsnvmereceiver

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

type Config struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`
}

var _ component.Config = (*Config)(nil)

func (c *Config) Validate() error {
	return nil
}
