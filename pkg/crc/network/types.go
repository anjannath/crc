package network

import (
	"fmt"

	"github.com/crc-org/crc/v2/pkg/crc/logging"
	"github.com/spf13/cast"
)

type NameServer struct {
	IPAddress string
}

type SearchDomain struct {
	Domain string
}

type ResolvFileValues struct {
	SearchDomains []SearchDomain
	NameServers   []NameServer
}

func (vals *ResolvFileValues) GetNameServer() []string {
	var nameservers []string
	for _, ns := range vals.NameServers {
		nameservers = append(nameservers, ns.IPAddress)
	}
	return nameservers
}

func (vals *ResolvFileValues) GetSearchDomains() []string {
	var searchDomains []string
	for _, sd := range vals.SearchDomains {
		searchDomains = append(searchDomains, sd.Domain)
	}
	return searchDomains
}

type Mode string

func (m Mode) String() string {
	return string(m)
}

const (
	SystemNetworkingMode Mode = "system"
	UserNetworkingMode   Mode = "user"
)

func parseMode(input string) (Mode, error) {
	switch input {
	case string(UserNetworkingMode), "vsock":
		return UserNetworkingMode, nil
	case string(SystemNetworkingMode), "default":
		return SystemNetworkingMode, nil
	default:
		return UserNetworkingMode, fmt.Errorf("Cannot parse mode '%s'", input)
	}
}
func ParseMode(input string) Mode {
	mode, err := parseMode(input)
	if err != nil {
		logging.Errorf("unexpected network mode %s, using default", input)
		mode = getDefaultMode()
	}
	return mode
}

func getDefaultMode() Mode {
	return UserNetworkingMode
}

func ValidateMode(val interface{}) (bool, string) {
	_, err := parseMode(cast.ToString(val))
	if err != nil {
		return false, fmt.Sprintf("network mode should be either %s or %s", SystemNetworkingMode, UserNetworkingMode)
	}
	return true, ""
}

func SuccessfullyAppliedMode(_ string, _ interface{}) string {
	return "Network mode changed. Please run `crc cleanup` and `crc setup`."
}
