package customers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// list customers
	getCustomersHandler := httptransport.NewServer(
		makeGetCustomersEndpoint(s),
		decodeGetCustomersRequest,
		encodeResponse,
		opts...,
	)

	// create customer
	createCustomerHandler := httptransport.NewServer(
		makeCreateCustomerEndpoint(s),
		decodeCreateCustomerRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/api/v1/customers", getCustomersHandler).Methods("GET")
	r.Handle("/api/v1/customers", createCustomerHandler).Methods("POST")

	return r
}

func decodeGetCustomersRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return getCustomersRequest{}, nil
}

func decodeCreateCustomerRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		return nil, err
	}
	return createCustomerRequest{Customer: c}, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func codeFrom(err error) int {
	switch err {
	case ErrInvalid:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
