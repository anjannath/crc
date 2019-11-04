package daemon

import (
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
	http.ListenAndServe("localhost:5732", r)
}
