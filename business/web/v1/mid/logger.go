package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/zmoog/shortener/foundation/web"
	"go.uber.org/zap"
)

// Logger writes information about the request to the logs.
func Logger(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v := web.GetValues(ctx)

			// Do some logging.
			log.Infow("request started", "trace_id", v.TraceID, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			// Do some more logging.
			log.Infow("request completed", "trace_id", v.TraceID, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr, "statuscode", v.StatusCode, "since", time.Since(v.Now))

			return err
		}

		return h
	}

	return m
}
