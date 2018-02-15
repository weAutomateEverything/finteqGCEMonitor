package gceservices

import (
	"github.com/go-kit/kit/metrics"
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

func (s *instrumentingService) RunServiceCheck(inward bool) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "RunServiceCheck").Add(1)
		s.requestLatency.With("method", "RunServiceCheck").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.RunServiceCheck(inward)
}
