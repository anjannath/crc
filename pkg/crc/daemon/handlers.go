package daemon

import (
	"fmt"
	"io/ioutil"

	"github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/crc/machine"
	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/cmd/crc/cmd/config"
	crcConfig "github.com/code-ready/crc/pkg/crc/config"
	"github.com/code-ready/crc/pkg/crc/validation"
	"github.com/code-ready/crc/pkg/crc/version"
)

func statusHandler() machine.ClusterStatusResult {
	statusConfig := machine.ClusterStatusConfig{Name: constants.DefaultName}
	clusterStatus, err := machine.Status(statusConfig)
	if err != nil {
		logging.Error(err.Error())
	}
	return clusterStatus
}

func stopHandler() machine.StopResult {
	stopConfig := machine.StopConfig{
		Name:  constants.DefaultName,
		Debug: true,
	}
	commandResult, err := machine.Stop(stopConfig)
	if err != nil {
		logging.Error(err.Error())
	}
	return commandResult
}

func startHandler() (machine.StartResult, error) {
	startConfig := machine.StartConfig{
		Name:          constants.DefaultName,
		BundlePath:    crcConfig.GetString(config.Bundle.Name),
		VMDriver:      crcConfig.GetString(config.VMDriver.Name),
		Memory:        crcConfig.GetInt(config.Memory.Name),
		CPUs:          crcConfig.GetInt(config.CPUs.Name),
		NameServer:    crcConfig.GetString(config.NameServer.Name),
		GetPullSecret: getPullSecretFileContent,
		Debug:         true,
	}
	status, err := machine.Start(startConfig)
	if err != nil {
		logging.Error(err.Error())
		return machine.StartResult{}, err
	}
	return status, nil
}

type versionResult struct {
	crcVersion string
	openshiftVersion string
}


func versionHandler() versionResult {
	var embedded string
	if !constants.BundleEmbedded() {
		embedded = "not "
	}
	return versionResult {
		crcVersion: fmt.Sprintf("crc version: %s+%s\n", version.GetCRCVersion(), version.GetCommitSha()),
		openshiftVersion: fmt.Sprintf("OpenShift version: %s (%sembedded in binary)\n", version.GetBundleVersion(), embedded),
	}
}

func getPullSecretFileContent() (string, error) {
	data, err := ioutil.ReadFile(crcConfig.GetString(config.PullSecretFile.Name))
	if err != nil {
		return "", err
	}
	pullsecret := string(data)
	if err := validation.ImagePullSecret(pullsecret); err != nil {
		return "", err
	}
	return pullsecret, nil
}
