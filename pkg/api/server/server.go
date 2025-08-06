package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jeremy2566/octopipe/internal/router"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type Config struct {
	Profile               string        `mapstructure:profile`
	Host                  string        `mapstructure:"host"`
	Port                  string        `mapstructure:"port"`
	PortMetrics           int           `mapstructure:"port-metrics"`
	ServerShutdownTimeout time.Duration `mapstructure:"server-shutdown-timeout"`
	HttpServerTimeout     time.Duration `mapstructure:"http-server-timeout"`
}

type Server struct {
	logger  *zap.Logger
	config  *Config
	handler http.Handler
}

func New(config *Config, logger *zap.Logger) *Server {
	srv := &Server{
		logger: logger,
		config: config,
	}

	return srv
}

func (s *Server) ListenAndServe() *http.Server {
	s.registerRoutes()
	// create the http server
	srv := s.startServer()
	return srv
}

func (s *Server) registerRoutes() {
	r := router.New(s.logger)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "health check is ok",
		})
		return
	})
	r.GET("/readyz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ready check is ok",
		})
	})

	s.handler = r
}

func (s *Server) startServer() *http.Server {
	srv := &http.Server{
		Addr:         s.config.Host + ":" + s.config.Port,
		WriteTimeout: s.config.HttpServerTimeout,
		ReadTimeout:  s.config.HttpServerTimeout,
		IdleTimeout:  2 * s.config.HttpServerTimeout,
		Handler:      s.handler,
	}

	go func() {
		s.logger.Info("Starting HTTP server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("HTTP server crashed", zap.Error(err))
		}
	}()
	return srv
}
