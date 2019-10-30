package account

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/shopspring/decimal"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

// New is logging wrapper for new account creation.
func (s *loggingService) New(id ID, currency Currency, balance decimal.Decimal) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "new",
			"id", id,
			"currency", currency,
			"balance", balance,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.New(id, currency, balance)
}

// Load is logging wrapper for load account.
func (s *loggingService) Load(id ID) (a *Account, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "load",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Load(id)
}

// LoadAll is logging wrapper for load all accounts.
func (s *loggingService) LoadAll() (r []*Account) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "loadAll",
			"len", len(r),
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.LoadAll()
}

// Delete is logging wrapper for delete account (mark it deleted).
func (s *loggingService) Delete(id ID) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "delete",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(id)
}
