package daemon

import (
	"net"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/handlers"
	"github.com/gorilla/rpc/json"
	"github.com/justinas/alice"

	"github.com/code-ready/crc/pkg/crc/logging"
)

func RunDaemon() {
	r := mux.NewRouter()
	// Create the server and register the service
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(CrcDaemon), "")

	chain := alice.New(
		func(h http.Handler) http.Handler {
			return handlers.CombinedLoggingHandler(logging.LogWriter(), h)
		})

	r.Handle("/rpc", chain.Then(s))
	os.RemoveAll("/tmp/crc.sock")
	lnr, err := net.Listen("unix", "/tmp/crc.sock")
	if err != nil {
		logging.Error("Failed to create socket")
	}
	http.Serve(lnr, r)
}
