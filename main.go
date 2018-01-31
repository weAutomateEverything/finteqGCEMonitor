package main

import (
	"github.com/zamedic/go2hal/alert"
	"os"
	"github.com/zamedic/go2hal/database"
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/cutofftimes"
	monitor2 "github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/monitor"
	"net/http"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os/signal"
	"syscall"
	"fmt"
)

func main() {

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	db := database.NewConnection()

	cutoffStore := cutofftimes.NewMongoStore(db)
	cutoffService := cutofftimes.NewService(cutoffStore, nil, true)

	alert := alert.NewKubernetesAlertProxy(errorEndpoint())

	_ = monitor2.NewService(alert, cutoffStore, seleniumServer())

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

func seleniumServer() string {
	return os.Getenv("SELENIUM_SERVER")
}

func errorEndpoint() string {
	return os.Getenv("HAL_ENDPOINT")
}
