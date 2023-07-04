package machine

import (
	"fmt"
	"path/filepath"

	"github.com/crc-org/crc/pkg/crc/constants"
	"github.com/crc-org/crc/pkg/crc/preset"
	crcos "github.com/crc-org/crc/pkg/os"
)

func copyDiskImage(destDir string, preset preset.Preset) (string, string, error) {
	const destFormat = "qcow2"

	imageName := fmt.Sprintf("%s.qcow2", constants.GetInstanceName(preset))

	srcPath := filepath.Join(constants.MachineInstanceDir, constants.GetInstanceName(preset), imageName)
	destPath := filepath.Join(destDir, imageName)

	_, _, err := crcos.RunWithDefaultLocale("qemu-img", "convert", "-f", "qcow2", "-O", destFormat, srcPath, destPath)
	if err != nil {
		return "", "", err
	}

	return destPath, destFormat, nil
}
