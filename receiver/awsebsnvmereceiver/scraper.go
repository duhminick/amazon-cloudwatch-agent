// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package awsebsnvmereceiver

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/aws/amazon-cloudwatch-agent/receiver/awsebsnvmereceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

const (
	nvmeDevicePrefix     = "nvme"
	devDirectoryPath     = "/dev"
	nvmeSysDirectoryPath = "/sys/class/nvme"
)

type nvmeScraper struct {
	logger *zap.Logger
	// mb *metadata.MetricsBuilder
}

func (s *nvmeScraper) start(_ context.Context, _ component.Host) error {
	s.logger.Debug("Starting NVME scraper", zap.String("receiver", metadata.Type.String()))
	return nil
}

func (s *nvmeScraper) shutdown(_ context.Context) error {
	s.logger.Debug("Shutting down NVME scraper", zap.String("receiver", metadata.Type.String()))
	return nil
}

func (s *nvmeScraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	s.logger.Info("[DOMINIC] Running NVME scraper")

	metric := pmetric.NewMetrics()
	// now := pcommon.NewTimestampFromTime(time.Now())

	return metric, nil
}

func getNvmeDevices() ([]string, error) {
	entries, err := os.ReadDir("/dev")
	if err != nil {
		return nil, err
	}

	devices := []string{}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), nvmeDevicePrefix) {
			devices = append(devices, entry.Name())
		}
	}

	return devices, nil
}

func getNvmeDeviceSerial(device string) (string, error) {
	data, err := os.ReadFile("/sys/class/nvme/" + device + "/serial")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Use mountinfo to get NVMe devices
// func getNvmeDevices() ([]string, error) {
// 	devices := []string{}
//
// 	infos, err := mountinfo.GetMounts(sourceFilter(nvmeDevicePrefix))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	for _, info := range infos {
// 		devices = append(devices, info.Source)
// 	}
//
// 	return devices, nil
// }
// func sourceFilter(prefix string) mountinfo.FilterFunc {
// 	return func(m *mountinfo.Info) (bool, bool) {
// 		skip := !strings.HasPrefix(m.Source, prefix)
// 		return skip, false
// 	}
// }

// Another option is to read `/sys/class/nvme/nvme{id}`. Inside has a serial file
// which will have the volume ID
// func getNvmeDevices() ([]string, error) {
// }

func newScraper(logger *zap.Logger) *nvmeScraper {
	return &nvmeScraper{
		logger: logger,
		// mb: metadata.NewMetricsBuilder()
	}
}
