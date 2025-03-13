package nvme

import (
	"fmt"
	"os"
	"strings"
)

func GetAllDevices() ([]NvmeDeviceFileAttributes, error) {
	entries, err := os.ReadDir(DevDirectoryPath)
	if err != nil {
		return nil, err
	}

	devices := []NvmeDeviceFileAttributes{}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), NvmeDevicePrefix) {
			device, err := ParseNvmeDeviceFileName(entry.Name())
			if err == nil {
				devices = append(devices, device)
			}
		}
	}

	return devices, nil
}

func GetDeviceSerial(device *NvmeDeviceFileAttributes) (string, error) {
	deviceName, err := device.BaseDeviceName()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(fmt.Sprintf("%s/%s/serial", NvmeSysDirectoryPath, deviceName))
	if err != nil {
		return "", err
	}
	return cleanupString(string(data)), nil
}

func GetDeviceModel(device *NvmeDeviceFileAttributes) (string, error) {
	deviceName, err := device.BaseDeviceName()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(fmt.Sprintf("%s/%s/model", NvmeSysDirectoryPath, deviceName))
	if err != nil {
		return "", err
	}
	return cleanupString(string(data)), nil
}

func IsEbsDevice(device *NvmeDeviceFileAttributes) (bool, error) {
	model, err := GetDeviceModel(device)
	if err != nil {
		return false, err
	}
	return model == EbsNvmeModelName, nil
}

func cleanupString(input string) string {
	// Some device info strings use fixed-width padding and/or end with a new line
	return strings.TrimSpace(strings.TrimSuffix(input, "\n"))
}
