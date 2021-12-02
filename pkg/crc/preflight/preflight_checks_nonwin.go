// +build !windows

package preflight

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/code-ready/crc/pkg/crc/logging"
	crcos "github.com/code-ready/crc/pkg/os"
)

var nonWinPreflightChecks = []Check{
	{
		configKeySuffix:  "check-root-user",
		checkDescription: "Checking if running as non-root",
		check:            checkIfRunningAsNormalUser,
		fixDescription:   "crc should not be run as root",
		flags:            NoFix,

		// no need for an "os" label as this is only built on relevant OSes through the use of golang build tags
		labels: None,
	},
}

func checkIfRunningAsNormalUser() error {
	if os.Geteuid() != 0 {
		return nil
	}
	logging.Debug("Ran as root")
	return errors.New("crc should not be run as root")
}

func setSuid(path string) error {
	logging.Debugf("Making %s suid", path)

	stdOut, stdErr, err := crcos.RunPrivileged(fmt.Sprintf("Changing ownership of %s", path), "chown", "root", path)
	if err != nil {
		return fmt.Errorf("Unable to set ownership of %s to root: %s %v: %s",
			path, stdOut, err, stdErr)
	}

	/* Can't do this before the chown as the chown will reset the suid bit */
	stdOut, stdErr, err = crcos.RunPrivileged(fmt.Sprintf("Setting suid for %s", path), "chmod", "u+s,g+x", path)
	if err != nil {
		return fmt.Errorf("Unable to set suid bit on %s: %s %v: %s", path, stdOut, err, stdErr)
	}
	return nil
}

func checkSuid(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.Mode()&os.ModeSetuid == 0 {
		return fmt.Errorf("%s does not have the SUID bit set (%s)", path, fi.Mode().String())
	}
	if fi.Sys().(*syscall.Stat_t).Uid != 0 {
		return fmt.Errorf("%s is not owned by root", path)
	}

	return nil
}

// Check if helper executable is cached or not
func checkAdminHelperExecutableCached() error {
	if version.IsInstaller() {
		return nil
	}

	helper := cache.NewAdminHelperCache()
	if !helper.IsCached() {
		return errors.New("crc-admin-helper executable is not cached")
	}
	if err := helper.CheckVersion(); err != nil {
		return errors.Wrap(err, "unexpected version of the crc-admin-helper executable")
	}
	logging.Debug("crc-admin-helper executable already cached")
	return checkSuid(helper.GetExecutablePath())
}

func fixAdminHelperExecutableCached() error {
	if version.IsInstaller() {
		return nil
	}

	helper := cache.NewAdminHelperCache()
	if err := helper.EnsureIsCached(); err != nil {
		return errors.Wrap(err, "Unable to download crc-admin-helper executable")
	}
	logging.Debug("crc-admin-helper executable cached")
	return setSuid(helper.GetExecutablePath())
}

var oldAdminHelpers = []string{"admin-helper-linux", "admin-helper-darwin"}

/* These 2 checks can be removed after a few releases */
func checkOldAdminHelperExecutableCached() error {
	logging.Debugf("Checking if an older admin-helper executable is installed")
	for _, oldExecutable := range oldAdminHelpers {
		oldPath := filepath.Join(constants.BinDir(), oldExecutable)
		if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
			return fmt.Errorf("Found old admin-helper executable '%s'", oldExecutable)
		}
	}

	logging.Debugf("No older admin-helper executable found")

	return nil
}

func fixOldAdminHelperExecutableCached() error {
	logging.Debugf("Removing older admin-helper executable")
	for _, oldExecutable := range oldAdminHelpers {
		oldPath := filepath.Join(constants.BinDir(), oldExecutable)
		if err := os.Remove(oldPath); err != nil {
			if !os.IsNotExist(err) {
				logging.Debugf("Failed to remove  %s: %v", oldPath, err)
				return err
			}
		} else {
			logging.Debugf("Successfully removed %s", oldPath)
		}
	}

	return nil
}

func checkSupportedCPUArch() error {
	if runtime.GOARCH != "amd64" {
		logging.Debugf("GOARCH is %s", runtime.GOARCH)
		return fmt.Errorf("CodeReady Containers can only run on x86_64 CPUs")
	}
	return nil
}
