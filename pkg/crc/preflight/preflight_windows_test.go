package preflight

import (
	"testing"

	"github.com/code-ready/crc/pkg/crc/config"
	"github.com/code-ready/crc/pkg/crc/network"
	"github.com/stretchr/testify/assert"
)

func TestCountConfigurationOptions(t *testing.T) {
	cfg := config.New(config.NewEmptyInMemoryStorage())
	RegisterSettings(cfg)
	assert.Len(t, cfg.AllConfigs(), 15)
}

func TestCountPreflights(t *testing.T) {
	assert.Len(t, getPreflightChecks(false, false, network.SystemNetworkingMode, []string{}), 20)
	assert.Len(t, getPreflightChecks(true, true, network.SystemNetworkingMode, []string{}), 20)
	assert.Len(t, getPreflightChecks(true, true, network.SystemNetworkingMode, []string{"add-user-to-hyperv-admin-group"}), 1)

	assert.Len(t, getPreflightChecks(false, false, network.UserNetworkingMode, []string{}), 21)
	assert.Len(t, getPreflightChecks(true, true, network.UserNetworkingMode, []string{}), 21)
}
