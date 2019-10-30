// Package db provide PostgreSQL storage for application data repositories.
// docker run -d -it -p 5432:5432 -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=db postgres:9.6
package db

import (
	"awesomeProject1/account"
	"awesomeProject1/errs"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// CreateSchema creating schema if its not exist. Without any migrations mechanic, just schema only.
func CreateSchema(conn *pg.DB) error {
	for _, model := range []interface{}{(*account.Account)(nil)} {
		err := conn.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type accountRepository struct {
	conn *pg.DB
}

// Store account in the repository
func (r *accountRepository) Store(account *account.Account) error {
	if err := r.conn.Insert(account); err != nil {
		return err
	}
	return nil
}

// Find account in the repository with specified id
func (r *accountRepository) Find(id account.ID) (*account.Account, error) {
	a := &account.Account{ID: id}
	err := r.conn.Select(a)
	if err != nil {
		return nil, err
	}
	if a.Deleted {
		return nil, errs.ErrUnknownAccount
	}
	return a, nil
}

// FindAll returns all accounts registered in the system
func (r *accountRepository) FindAll() []*account.Account {
	var accounts []*account.Account
	err := r.conn.Model(&accounts).Where("deleted = ?", false).Select()
	if err != nil {
		return nil
	}
	return accounts
}

// MarkDeleted is mark as deleted specified account in the system
func (r *accountRepository) MarkDeleted(id account.ID) error {
	a := &account.Account{ID: id}
	err := r.conn.Select(a)
	if err != nil {
		return err
	}
	if a.ID == "" || a.Deleted {
		return errs.ErrUnknownAccount
	}
	a.Deleted = true
	err = r.conn.Update(a)
	if err != nil {
		return err
	}
	return nil
}

// NewAccountRepository returns a new instance of a PostgreSQL account repository.
func NewAccountRepository(conn *pg.DB) account.Repository {
	return &accountRepository{
		conn: conn,
	}
}
