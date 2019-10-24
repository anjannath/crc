package daemon

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

	"github.com/code-ready/crc/pkg/crc/logging"
)

func RunDaemon() {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	// Services offered by the daemon
	status := new(ClusterStatus)
	stop := new(ClusterStop)
	start := new(ClusterStart)
	s.RegisterService(status, "")
	s.RegisterService(stop, "")
	s.RegisterService(start, "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	logging.Info("============Launched server==========")
	http.ListenAndServe(":1234", r)
}
