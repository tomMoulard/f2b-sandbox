package fail2ban

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/http-wasm/http-wasm-guest-tinygo/handler"
	"github.com/http-wasm/http-wasm-guest-tinygo/handler/api"
)

type Middleware struct {
	Config     *Config
	middleware *Fail2Ban
}

var mw = &Middleware{}

func init() {
	err := json.Unmarshal(handler.Host.GetConfig(), &mw.Config)
	if err != nil {
		handler.Host.Log(api.LogLevelError, fmt.Sprintf("Could not load config %v", err))
		os.Exit(1)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	hander, err := New(context.Background(), next, mw.Config, "fail2ban-WASM")
	if err != nil {
		handler.Host.Log(api.LogLevelError, fmt.Sprintf("Could create middleware: %v", err))
		os.Exit(1)
	}

	var ok bool
	mw.middleware, ok = hander.(*Fail2Ban)
	if !ok {
		handler.Host.Log(api.LogLevelError, "Could create middleware")
		os.Exit(1)
	}
}

func main() {
	handler.HandleRequestFn = mw.handleRequest
	// handler.HandleResponseFn = mw.handleResponse
}

// handleRequest implements a simple request middleware.
// Wraps the Fail2ban plugin.
func (mw *Middleware) handleRequest(req api.Request, resp api.Response) (next bool, reqCtx uint32) {
	next = mw.middleware.shouldAllow(req.GetSourceAddr(), req.GetURI())
	return
}

// handleResponse implements a simple response middleware.
// NOOP for this particular plugin.
// func (mw *Middleware) handleResponse(_ uint32, _ api.Request, _ api.Response, _ bool) {
// }
