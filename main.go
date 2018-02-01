package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/frankgreco/demogods/pkg/log"
	"github.com/frankgreco/demogods/pkg/metrics"
	"github.com/frankgreco/demogods/pkg/middleware"
	"github.com/frankgreco/demogods/pkg/server"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithContext(r.Context())
	if rand.Intn(100) < 20 {
		metrics.Demos.WithLabelValues("fail").Inc()
		logger.Warn("demo failed")
		w.Write([]byte("You live demo will fail!"))
	} else {
		metrics.Demos.WithLabelValues("succeed").Inc()
		logger.Info("demo succeeded")
		w.Write([]byte("You live demo will succeed!"))
	}
}

func logError(err error) error {
	if err != nil {
		log.WithContext(nil).Error(err.Error())
	}
	return err
}

func main() {
	ctx, cancel := context.WithCancel(server.SetupSignalHandler())

	demogods := server.Prepare(ctx, &server.Options{
		Name:    "demogods",
		Addr:    "0.0.0.0",
		Port:    "8080",
		Handler: middleware.Correlation(middleware.Metrics(http.HandlerFunc(myHandler))),
	})

	metrics := server.Prepare(ctx, &server.Options{
		Name:    "prometheus",
		Addr:    "0.0.0.0",
		Port:    "9000",
		Handler: promhttp.Handler(),
	})

	var g run.Group

	g.Add(func() error {
		<-ctx.Done()
		return nil
	}, func(error) {
		cancel()
	})

	g.Add(func() error {
		return logError(metrics.Run())
	}, func(error) {
		logError(metrics.Close())
	})

	g.Add(func() error {
		return logError(demogods.Run())
	}, func(error) {
		logError(demogods.Close())
	})

	if err := g.Run(); err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
