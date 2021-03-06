package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Masterminds/semver/v3"
)

// The following variables are private fields and should be set when compiling with ldflags, for example --ldflags="-X github.com/code-ready/crc/pkg/version.crcVersion=vX.Y.Z
var (
	// The current version of minishift
	crcVersion = "0.0.0-unset"

	// The SHA-1 of the commit this executable is build off
	commitSha = "sha-unset"

	// Bundle version which used for the release.
	bundleVersion = "0.0.0-unset"

	okdBuild = "false"

	macosInstallPath = "/unset"

	msiBuild = "false"
)

const (
	releaseInfoLink = "https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/crc/latest/release-info.json"
	// Tray version to be embedded in executable
	crcMacTrayVersion = "1.0.9"
	// Windows forms application version type major.minor.buildnumber.revesion
	crcWindowsTrayVersion = "0.7.0.0"
)

type CrcReleaseInfo struct {
	Version Version           `json:"version"`
	Links   map[string]string `json:"links"`
}

type Version struct {
	CrcVersion       *semver.Version `json:"crcVersion"`
	GitSha           string          `json:"gitSha"`
	OpenshiftVersion string          `json:"openshiftVersion"`
}

func GetCRCVersion() string {
	return crcVersion
}

func GetCommitSha() string {
	return commitSha
}

func GetBundleVersion() string {
	return bundleVersion
}

func IsOkdBuild() bool {
	return okdBuild == "true"
}

func GetCRCMacTrayVersion() string {
	return crcMacTrayVersion
}

func GetCRCWindowsTrayVersion() string {
	return crcWindowsTrayVersion
}

func GetMacosInstallPath() string {
	return macosInstallPath
}

func IsMacosInstallPathSet() bool {
	return macosInstallPath != "/unset"
}

func IsMsiBuild() bool {
	return msiBuild != "false"
}

func GetCRCLatestVersionFromMirror() (*CrcReleaseInfo, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(releaseInfoLink)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP error: %s: %d", response.Status, response.StatusCode)
	}

	releaseMetaData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var releaseInfo CrcReleaseInfo
	if err := json.Unmarshal(releaseMetaData, &releaseInfo); err != nil {
		return nil, fmt.Errorf("Error unmarshaling JSON metadata: %v", err)
	}

	return &releaseInfo, nil
}
