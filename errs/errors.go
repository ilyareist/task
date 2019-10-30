// Package errs contain types and methods to help in error-handling.
package errs

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrUnknownAccount       = errors.New("unknown account")
	ErrInvalidArgument      = errors.New("invalid argument")
	ErrUnknownSourceAccount = errors.New("unknown source account")
	ErrUnknownTargetAccount = errors.New("unknown target account")
	ErrAccountsAreEqual     = errors.New("target account must not be equal to source account")
	ErrInsufficientMoney    = errors.New("insufficient money on source account")
	ErrStorePayments        = errors.New("can not store payments")
	ErrStoreSourceAccount   = errors.New("can not update source account")
	ErrStoreTargetAccount   = errors.New("can not update target account")
	ErrBadRoute             = errors.New("bad route")
)

// ValidationError represents validation error, for right choosing of HTTP status in response.
type ValidationError struct {
	Err error
}

// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
func (e ValidationError) Error() string {
	return "validation error: " + e.Err.Error()
}

type errorer interface {
	ErrError() error
}

// EncodeResponse encoding response data by default way (without any struct). It is enough in most cases.
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	//if e, ok := response.(errorer); ok && e.error() != nil {
	e, ok := response.(errorer)
	if ok && e.ErrError() != nil {
		EncodeError(ctx, e.ErrError(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// EncodeError encode errs from business-logic
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrUnknownAccount, ErrUnknownSourceAccount, ErrUnknownTargetAccount:
		w.WriteHeader(http.StatusNotFound)
	case ErrInvalidArgument, ErrInsufficientMoney:
		w.WriteHeader(http.StatusBadRequest)
	case ErrAccountsAreEqual:
		w.WriteHeader(http.StatusNotAcceptable)
	default:
		switch err.(type) {
		case ValidationError:
			w.WriteHeader(http.StatusNotAcceptable)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// ErrorOnlyResponse represents response which may contain only error or nothing.
type ErrorOnlyResponse struct {
	Err error `json:"error,omitempty"`
}

func (r ErrorOnlyResponse) error() error { return r.Err }
