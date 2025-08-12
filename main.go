package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jeremy2566/octopipe/internal/service"
	"github.com/jeremy2566/octopipe/pkg/api/server"
	"github.com/jeremy2566/octopipe/pkg/logger"
	"github.com/jeremy2566/octopipe/pkg/shutdown"
	"github.com/jeremy2566/octopipe/pkg/signals"
	"github.com/jeremy2566/octopipe/pkg/version"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"resty.dev/v3"
)

func main() {
	fs := pflag.NewFlagSet("default", pflag.ContinueOnError)
	fs.String("profile", "dev", "dev, test, uat or prod")
	fs.String("host", "127.0.0.1", "HTTP host to bind service to")
	fs.Int("port", 6652, "HTTP port to bind service to")
	fs.Int("port-metrics", 0, "metrics port")
	fs.Duration("server-shutdown-timeout", 5*time.Second, "server graceful shutdown timeout duration")
	fs.Duration("http-server-timeout", 30*time.Second, "server read and write timeout duration")

	fs.String("level", "info", "log level debug, info, warn, error, fatal or panic")
	fs.String("config-path", "", "config dir path")
	fs.String("config", "config.yaml", "config file name")

	vFLag := fs.BoolP("version", "v", false, "show version and exit")
	err := fs.Parse(os.Args[1:])
	switch {
	case errors.Is(err, pflag.ErrHelp):
		os.Exit(0)
	case err != nil:
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		fs.PrintDefaults()
		os.Exit(2)
	case *vFLag:
		fmt.Println(version.VERSION)
		os.Exit(0)
	}
	err = viper.BindPFlags(fs)
	viper.Set("version", version.VERSION)
	viper.Set("revision", version.REVISION)
	viper.AutomaticEnv()

	// load config from file
	if _, fileErr := os.Stat(filepath.Join(viper.GetString("config-path"), viper.GetString("config"))); fileErr == nil {
		viper.SetConfigName(strings.Split(viper.GetString("config"), ".")[0])
		viper.AddConfigPath(viper.GetString("config-path"))
		if readErr := viper.ReadInConfig(); readErr != nil {
			fmt.Printf("Error reading config file, %v\n", readErr)
		}
	}

	log, _ := logger.Load(viper.GetString("level"))
	defer log.Sync()

	stdLog := zap.RedirectStdLog(log)
	defer stdLog()

	// validate port
	if _, err := strconv.Atoi(viper.GetString("port")); err != nil {
		port, _ := fs.GetInt("port")
		viper.Set("port", strconv.Itoa(port))
	}

	cfg := server.Config{
		Users: map[string]string{},
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panic("config unmarshal failed.", zap.Error(err))
	}

	// log version and port
	log.Info("Starting echo",
		zap.String("version", viper.GetString("version")),
		zap.String("revision", viper.GetString("revision")),
		zap.String("port", cfg.Port),
	)

	srv := server.New(&cfg, log)
	httpServer := srv.ListenAndServe()

	go func() {
		srv := service.NewCache(
			log,
			resty.New().SetRetryCount(3).SetRetryWaitTime(1*time.Second).SetRetryMaxWaitTime(5*time.Second),
		)
		err := srv.SyncCache("fat-base-envrionment")
		if err != nil {
			log.Error("sync cache failed.", zap.Error(err))
		}
	}()

	stopCh := signals.SetupSignalHandler()
	sd, _ := shutdown.New(cfg.ServerShutdownTimeout, log)
	sd.Graceful(stopCh, httpServer)
}
