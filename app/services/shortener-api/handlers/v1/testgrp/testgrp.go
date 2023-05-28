package testgrp

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	v1Web "github.com/zmoog/shortener/business/web/v1"
	"github.com/zmoog/shortener/foundation/web"
)

// Status represents a test handler for now.
func Status(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return v1Web.NewRequestError(fmt.Errorf("trusted error"), http.StatusBadRequest)
		// panic("oh no we panicked")
	}

	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
