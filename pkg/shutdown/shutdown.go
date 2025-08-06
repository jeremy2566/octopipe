package shutdown

import (
	"context"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Shutdown struct {
	logger                *zap.Logger
	serverShutdownTimeout time.Duration
}

func New(serverShutdownTimeout time.Duration, logger *zap.Logger) (*Shutdown, error) {
	srv := &Shutdown{
		logger:                logger,
		serverShutdownTimeout: serverShutdownTimeout,
	}

	return srv, nil
}

func (s *Shutdown) Graceful(
	stopCh <-chan struct{},
	httpServer *http.Server,
) {
	ctx := context.Background()

	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(ctx, s.serverShutdownTimeout)
	defer cancel()

	s.logger.Info("Shutting down HTTP/HTTPS server", zap.Duration("timeout", s.serverShutdownTimeout))

	if viper.GetString("level") != "debug" {
		time.Sleep(3 * time.Second)
	}

	// determine if the http server was started
	if httpServer != nil {
		if err := httpServer.Shutdown(ctx); err != nil {
			s.logger.Warn("HTTP server graceful shutdown failed", zap.Error(err))
		}
	}
}
