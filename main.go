package main

import (
	"fmt"
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/cutofftimes"
	monitor2 "github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/monitor"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/zamedic/go2hal/alert"
	"github.com/zamedic/go2hal/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/gceSelenium"
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/gceservices"
)

func main() {

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	db := database.NewConnection()

	fieldKeys := []string{"method"}

	alert := alert.NewKubernetesAlertProxy(errorEndpoint())

	seleniumService := gceSelenium.NewService(alert)
	seleniumService = gceSelenium.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "selenium",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "selenium",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "selenium",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), seleniumService)

	gceService := gceservices.NewService(seleniumService)
	gceService = gceservices.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "gceService",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "gceService",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "gceService",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), gceService)

	cutoffStore := cutofftimes.NewMongoStore(db)
	cutoffService := cutofftimes.NewService(cutoffStore, seleniumService)
	cutoffService = cutofftimes.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "cutoffService",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "cutoffService",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "cutoffService",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), cutoffService)

	_ = monitor2.NewService(alert, seleniumService, gceService, cutoffService)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/cutoff", cutofftimes.MakeHandler(cutoffService, httpLogger))

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", ":8000", "msg", "listening")
		errs <- http.ListenAndServe(":8000", nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)

}

func errorEndpoint() string {
	return os.Getenv("HAL_ENDPOINT")
}
