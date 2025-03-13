package nvme

import (
	"fmt"
	"os"
	"strings"
)

func GetAllDevices() ([]string, error) {
	entries, err := os.ReadDir(DevDirectoryPath)
	if err != nil {
		return nil, err
	}

	devices := []string{}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), NvmeDevicePrefix) {
			devices = append(devices, entry.Name())
		}
	}

	return devices, nil
}

func GetDeviceSerial(device string) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("%s/%s/serial", NvmeSysDirectoryPath, device))
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(data), "\n"), nil
}

func GetDeviceModel(device string) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("%s/%s/model", NvmeSysDirectoryPath, device))
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(data), "\n"), nil
}

func IsEbsDevice(device string) (bool, error) {
	model, err := GetDeviceModel(device)
	if err != nil {
		return false, err
	}
	return model == EbsNvmeModelName, nil
}
