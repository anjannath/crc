//go:build !linux
// +build !linux

package machine

import (
	"fmt"
	"runtime"

	"github.com/crc-org/crc/pkg/crc/preset"
)

func copyDiskImage(dirName string, _ preset.Preset) (string, string, error) {
	return "", "", fmt.Errorf("Not implemented for %s", runtime.GOOS)
}
