package account

import (
	"context"
	"encoding/json"
	"net/http"

	"../errs"
	"github.com/shopspring/decimal"

	"github.com/asaskevich/govalidator"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func init() {
	// How to serialize decimals to JSON
	decimal.MarshalJSONWithoutQuotes = true

	// Decimal validator plugin for govaidator
	govalidator.CustomTypeTagMap.Set("decimal", func(i interface{}, context interface{}) bool {
		switch i.(type) {
		case decimal.Decimal:
			return true
		}
		return false
	})
}

// MakeHandler returns a handler for the account service.
func MakeHandler(as Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(errs.EncodeError),
	}

	newAccountHandler := kithttp.NewServer(
		makeNewAccountEndpoint(as),
		decodeNewAccountRequest,
		errs.EncodeResponse,
		opts...,
	)

	loadAccountHandler := kithttp.NewServer(
		makeLoadAccountEndpoint(as),
		decodeLoadAccountRequest,
		errs.EncodeResponse,
		opts...,
	)

	loadAllAccountsHandler := kithttp.NewServer(
		makeLoadAllAccountsEndpoint(as),
		decodeLoadAllAccountsRequest,
		errs.EncodeResponse,
		opts...,
	)

	deleteAccountHandler := kithttp.NewServer(
		makeDeleteAccountEndpoint(as),
		decodeDeleteAccountRequest,
		errs.EncodeResponse,
		opts...,
	)

	router := mux.NewRouter()

	router.Handle("/api/accounts/v1/accounts", newAccountHandler).Methods("POST")
	router.Handle("/api/accounts/v1/accounts", loadAllAccountsHandler).Methods("GET")
	router.Handle("/api/accounts/v1/accounts/{id}", loadAccountHandler).Methods("GET")
	router.Handle("/api/accounts/v1/accounts/{id}", deleteAccountHandler).Methods("DELETE")

	return router
}

func decodeNewAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body newAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	if _, err := govalidator.ValidateStruct(body); err != nil {
		return nil, errs.ValidationError{Err: err}
	}
	return body, nil
}

func decodeLoadAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errs.ErrBadRoute
	}
	return idField{ID: ID(id)}, nil
}

func decodeLoadAllAccountsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeDeleteAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errs.ErrBadRoute
	}
	return idField{ID: ID(id)}, nil
}
