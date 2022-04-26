package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	redis "github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/vontikov/prom-redis/internal/collector"
	"github.com/vontikov/prom-redis/internal/conn"
	"github.com/vontikov/prom-redis/internal/logging"
	"github.com/vontikov/prom-redis/internal/parser"
	"github.com/vontikov/prom-redis/internal/util"
)

var (
	// App is the app name.
	App string = "N/A"
	// Version is the app version.
	Version string = "N/A"
)

var (
	infoSectionFlag   = flag.String("section", "everything", "Information and statistics section.")
	listenAddrFlag    = flag.String("listen-address", ":3501", "The address to listen on for HTTP requests.")
	logLevelFlag      = flag.String("log-level", "info", "Log level: trace|debug|info|warn|error|none")
	redisAddrFlag     = flag.String("address", "localhost:6379", "Redis address to connect, host:port.")
	redisPasswordFlag = flag.String("password", "", "Redis password.")
	redisUsernameFlag = flag.String("username", "", "Redis username.")
	metricNsFlag      = flag.String("namespace", "redis", "Metric namespace.")
)

func main() {
	flag.Parse()

	logging.SetLevel(*logLevelFlag)
	logger := logging.NewLogger(App)

	hostname, err := os.Hostname()
	util.PanicOnError(err)

	logger.Info("starting", "name", App, "version", Version, "hostname", hostname)
	flag.VisitAll(func(f *flag.Flag) { logger.Debug("option", "name", f.Name, "value", f.Value) })

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	opts := &redis.Options{
		Addr:     *redisAddrFlag,
		Username: *redisUsernameFlag,
		Password: *redisPasswordFlag,
		DB:       0,
	}
	rc := conn.New(ctx, opts, *infoSectionFlag)
	prometheus.MustRegister(
		collector.NewCollector(rc.Importer(), parser.Parse, *metricNsFlag))

	router := mux.NewRouter()
	router.PathPrefix("/metrics").Handler(promhttp.Handler())

	srv := &http.Server{
		Addr:    *listenAddrFlag,
		Handler: router,
	}

	logger.Info("listening", "address", *listenAddrFlag)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen: %s\n", err)
		}
		logger.Info("shutdown")
	}()

	await := make(chan any)
	go func() {
		<-ctx.Done()
		srv.Close()
		await <- 1
	}()

	sig := <-signals
	logger.Info("received", "signal", sig)
	cancel()
	<-await
}
