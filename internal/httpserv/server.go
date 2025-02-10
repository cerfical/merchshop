package httpserv

import (
	"context"
	"errors"
	stdlog "log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/cerfical/merchshop/internal/log"
)

func New(cfg *Config, h http.Handler, log *log.Logger) *Server {
	servAddr := net.JoinHostPort(cfg.Host, cfg.Port)
	return &Server{
		serv: http.Server{
			Addr: servAddr,

			// Log requests before any other logic applies
			Handler: logRequest(log)(h),

			// Redirect [http.Server] errors to a custom logger
			ErrorLog: stdlog.New(&httpErrorLog{log}, "", 0),
		},
		log: log,
	}
}

type httpErrorLog struct {
	log *log.Logger
}

func (w *httpErrorLog) Write(p []byte) (int, error) {
	// Trim carriage return produced by stdlog
	n := len(p)
	if n > 0 && p[n-1] == '\n' {
		p = p[0 : n-1]
		n--
	}

	w.log.Error("HTTP server failure", errors.New(string(p)))
	return n, nil
}

type Server struct {
	serv http.Server
	log  *log.Logger
}

func (s *Server) Run(ctx context.Context) (err error) {
	s.log.WithFields("addr", s.serv.Addr).Info("Starting up the server")

	sigCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error)
	go func() {
		if err := s.serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	defer func() {
		if closeErr := s.serv.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	select {
	case <-sigCtx.Done():
		// The server stopped due to a system signal, try to shutdown the server cleanly
		s.log.Info("Shutting down the server")
		if err := s.serv.Shutdown(ctx); err != nil {
			s.log.Error("Failed to shut down the server", err)
			return err
		}
	case err := <-errChan:
		// The server was terminated abnormally, exit now
		return err
	}

	return nil
}
