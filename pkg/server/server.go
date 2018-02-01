package server

import (
	"context"
	"net"
	"net/http"

	"go.uber.org/zap"

	"github.com/frankgreco/demogods/pkg/log"
)

// Options contains configuration for an HTTP server
type Options struct {
	Name    string
	Addr    string
	Port    string
	Handler http.Handler

	ctx    context.Context
	server *http.Server
	err    error
}

// Prepare will create a new server that is ready to be started.
// If an error is encountered, it will be abstracted by the
// returned type to be handled by future function calls.
func Prepare(ctx context.Context, opts *Options) *Options {
	opts.ctx = ctx
	opts.server = &http.Server{
		Addr:    net.JoinHostPort(opts.Addr, opts.Port),
		Handler: opts.Handler,
	}

	return opts
}

// Run will start an HTTP sever. If an error is encountered,
// it will be returned immedietaly.
func (opts *Options) Run() error {
	if opts.err != nil {
		return opts.err
	}

	log.WithContext(nil).Info("starting server",
		zap.String("name", opts.Name),
		zap.String("address", opts.server.Addr),
	)
	err := opts.server.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Close will gracefaully close the HTTP server. Any in flight requests
// will be given time to complete.
func (opts *Options) Close() error {
	log.WithContext(nil).Info("closing server",
		zap.String("name", opts.Name),
		zap.String("address", opts.server.Addr),
	)
	return opts.server.Shutdown(opts.ctx)
}
