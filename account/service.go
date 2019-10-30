// Package account provides handlers for work with accounts in the system.
package account

import (
	"github.com/shopspring/decimal"
)

// Currency type represents available currencies
type Currency string

const (
	CurrencyUSD Currency = "USD"
)

// ID type used for accounts identification.
type ID string

// Account is a wallet in the system.
type Account struct {
	TableName struct{}        `json:"-" sql:"select:accounts_view,alias:accounts"`
	ID        ID              `json:"id" sql:"id,pk,type:varchar(255)"`
	Balance   decimal.Decimal `json:"balance" sql:"balance,notnull,type:'decimal(16,4)'"`
	Currency  Currency        `json:"currency" sql:"currency,notnull,type:varchar(3)"`
	Deleted   bool            `json:"-" sql:"deleted,notnull"`
}

// Service is the interface that provides account methods.
type Service interface {
	// New registers a new account in the system, with desired Balance.
	New(id ID, currency Currency, balance decimal.Decimal) error

	// Load returns a read model of an account.
	Load(id ID) (*Account, error)

	// LoadAll returns all accounts registered in the system.
	LoadAll() []*Account

	// Delete uses to delete account from the system. Actually mark it as deleted.
	Delete(id ID) error
}

type service struct {
	accounts Repository
}

// New registers a new account in the system, with zero Balance.
func (s *service) New(id ID, currency Currency, balance decimal.Decimal) error {
	if currency == "" {
		currency = CurrencyUSD
	}
	return s.accounts.Store(&Account{
		ID:       id,
		Balance:  balance,
		Currency: currency,
	})
}

// Load returns a read model of an account.
func (s *service) Load(id ID) (*Account, error) {
	a, err := s.accounts.Find(id)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// LoadAll returns all accounts registered in the system.
func (s *service) LoadAll() []*Account {
	return s.accounts.FindAll()
}

// Delete uses to delete account from the system. Actually mark it as deleted.
func (s *service) Delete(id ID) error {
	return s.accounts.MarkDeleted(id)
}

// NewService creates an account service with necessary dependencies.
func NewService(accounts Repository) Service {
	return &service{
		accounts: accounts,
	}
}

// Repository interface for accounts storing and operations.
type Repository interface {
	// Store account in the repository
	Store(account *Account) error

	// Find account in the repository with specified id
	Find(id ID) (*Account, error)

	// FindAll returns all accounts registered in the system
	FindAll() []*Account

	// MarkDeleted is mark as deleted specified account in the system
	MarkDeleted(id ID) error
}
