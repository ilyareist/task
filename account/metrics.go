package account

import (
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/shopspring/decimal"
)

type metricsService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

// NewMetricsService returns an instance of a metrics Service.
func NewMetricsService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &metricsService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

// New is logging wrapper for new account creation.
func (s *metricsService) New(id ID, currency Currency, balance decimal.Decimal) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "new").Add(1)
		s.requestLatency.With("method", "new").Observe(time.Since(begin).Seconds() * 100000)
	}(time.Now())

	return s.Service.New(id, currency, balance)
}

// Load is logging wrapper for load account.
func (s *metricsService) Load(id ID) (*Account, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "load").Add(1)
		s.requestLatency.With("method", "load").Observe(time.Since(begin).Seconds() * 100000)
	}(time.Now())

	return s.Service.Load(id)
}

// LoadAll is logging wrapper for load all accounts.
func (s *metricsService) LoadAll() []*Account {
	defer func(begin time.Time) {
		s.requestCount.With("method", "loadAll").Add(1)
		s.requestLatency.With("method", "loadAll").Observe(time.Since(begin).Seconds() * 100000)
	}(time.Now())

	return s.Service.LoadAll()
}

// Delete is logging wrapper for delete account (mark it deleted).
func (s *metricsService) Delete(id ID) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "delete").Add(1)
		s.requestLatency.With("method", "delete").Observe(time.Since(begin).Seconds() * 100000)
	}(time.Now())

	return s.Service.Delete(id)
}
