package daemon

import (
	"fmt"
	"io/ioutil"
	"encoding/json"

	"github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/crc/machine"
	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/cmd/crc/cmd/config"
	crcConfig "github.com/code-ready/crc/pkg/crc/config"
	"github.com/code-ready/crc/pkg/crc/validation"
	"github.com/code-ready/crc/pkg/crc/version"
)

func statusHandler() string {
	statusConfig := machine.ClusterStatusConfig{Name: constants.DefaultName}
	clusterStatus, err := machine.Status(statusConfig)
	if err != nil {
		logging.Error(err.Error())
	}
	s, err := json.Marshal(clusterStatus)
	if err != nil {
		return "Failed"
	}
	return string(s)
}

func stopHandler() string {
	stopConfig := machine.StopConfig{
		Name:  constants.DefaultName,
		Debug: true,
	}
	commandResult, err := machine.Stop(stopConfig)
	if err != nil {
		logging.Error(err.Error())
	}
	s, err := json.Marshal(commandResult)
	if err != nil {
		return "Failed"
	}
	return string(s)
}

func startHandler() (string, error) {
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
	fmt.Println(crcConfig.GetString(config.PullSecretFile.Name))
	status, err := machine.Start(startConfig)
	if err != nil {
		logging.Error(err.Error())
		return "Failed", err
	}
	s, err := json.Marshal(status)
	if err != nil {
		return "Failed during json marshaling", err
	}
	return string(s), nil
}

type VersionResult struct {
	CrcVersion string
	OpenshiftVersion string
}


func versionHandler() string {
	var embedded string
	if !constants.BundleEmbedded() {
		embedded = "not "
	}
	v := &VersionResult {
		CrcVersion: fmt.Sprintf("%s+%s", version.GetCRCVersion(), version.GetCommitSha()),
		OpenshiftVersion: fmt.Sprintf("%s (%sembedded in binary)", version.GetBundleVersion(), embedded),
	}
	s, err := json.Marshal(v)
	if err != nil {
		return "Failed"
	}
	return string(s)
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
