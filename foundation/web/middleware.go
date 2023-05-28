package web

// Middleware is a function that wraps a handler to perform pre/post processing.
type Middleware func(Handler) Handler

// wrapMiddleware creates a new handler by wrapping middleware around a final handler.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}

	}

	return handler
}
