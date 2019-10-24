package daemon

import (
	"net/http"
	"github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/crc/machine"
)

type ClusterStatus struct {}

func (t *ClusterStatus) GetStatus(r *http.Request, Args *machine.ClusterStatusConfig, Response *machine.ClusterStatusResult) error {
	logging.Debugf("Args: %#v", Args)
	s, err := machine.Status(*Args)
	if err !=nil {
		logging.Error(err.Error())
		return err
	}
	*Response = s
	return nil
}

type ClusterStop struct {}

func (t *ClusterStop) PerformStop(r *http.Request, Args *machine.StopConfig, Response *machine.StopResult) error {
	logging.Debugf("Args: %#v", Args)
	s, err := machine.Stop(*Args)
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	*Response = s
	return nil
}

type ClusterStartConfig struct {
	Name string
	BundlePath string
	VMDriver string
	Memory   int
	CPUs     int
	NameServer string
	Debug bool
	PullSecret string
}

type ClusterStart struct {}

func (t *ClusterStart) PerformStart(r *http.Request, Args *ClusterStartConfig, Response *machine.StartResult) error {
	logging.Debugf("Args: %#v", Args)
	startConfig := machine.StartConfig{
		Name: Args.Name,
		BundlePath: Args.BundlePath,
		VMDriver: Args.VMDriver,
		Memory: Args.Memory,
		CPUs: Args.CPUs,
		NameServer: Args.NameServer,
		Debug: Args.Debug,
		GetPullSecret: func() (string, error) { return Args.PullSecret, nil },
	}
	logging.Debug()
	logging.Debugf("%#v", startConfig)
	status, err := machine.Start(startConfig)
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	*Response = status
	return nil
}
