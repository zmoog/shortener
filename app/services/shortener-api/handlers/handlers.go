// Package handlers manages the different versions of the API.
package handlers

import (
	"net/http"
	"os"

	"github.com/zmoog/shortener/app/services/shortener-api/handlers/v1/testgrp"
	"github.com/zmoog/shortener/business/web/v1/mid"
	"github.com/zmoog/shortener/foundation/web"
	"go.uber.org/zap"
)

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	app.Handle(http.MethodGet, "/status", testgrp.Status)

	return app
}
