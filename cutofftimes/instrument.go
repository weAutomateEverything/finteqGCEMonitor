package cutofftimes

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

func (s *instrumentingService) DoCheck(inward bool) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "DoCheck").Add(1)
		s.requestLatency.With("method", "DoCheck").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.DoCheck(inward)
}

func (s *instrumentingService) parseInwardCutttoffTimes(request string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "parseInwardCutttoffTimes").Add(1)
		s.requestLatency.With("method", "parseInwardCutttoffTimes").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.parseInwardCutttoffTimes(request)
}
