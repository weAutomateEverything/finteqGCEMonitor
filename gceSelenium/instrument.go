package gceSelenium

import (
	"github.com/go-kit/kit/metrics"
	"github.com/tebeka/selenium"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	errorCount     metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentService(counter metrics.Counter, errorCount metrics.Counter,
	latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		errorCount:     errorCount,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) HandleSeleniumError(internal bool, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "HandleSeleniumError").Add(1)
		s.requestLatency.With("method", "HandleSeleniumError").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.HandleSeleniumError(internal, err)
}

func (s *instrumentingService) NewClient() (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "DoCheck").Add(1)
		s.requestLatency.With("method", "DoCheck").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "DoCheck").Add(1)
		}
	}(time.Now())
	return s.Service.NewClient()
}

func (s *instrumentingService) Driver() selenium.WebDriver {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Driver").Add(1)
		s.requestLatency.With("method", "Driver").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Driver()
}

func (s *instrumentingService) WaitForWaitFor() (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "WaitForWaitFor").Add(1)
		s.requestLatency.With("method", "WaitForWaitFor").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "WaitForWaitFor").Add(1)
		}
	}(time.Now())
	return s.Service.WaitForWaitFor()
}
