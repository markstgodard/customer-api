package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/markstgodard/customer-api/customers"
)

const (
	// BANNER is what is printed for help/info output

	BANNER = `
CUSTOMER API

`
)

var (
	port int
)

func init() {
	fmt.Printf(BANNER)

	// parse flags
	flag.IntVar(&port, "port", 8080, "HTTP listen address")

	flag.Parse()
}

func main() {

	// logger
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	// create service
	repo := &customers.InMemoryCustomerRepo{}
	s := customers.NewService(repo)
	s = customers.NewLoggingMiddlware(logger)(s)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/api/v1/", customers.MakeHTTPHandler(s, httpLogger))

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		addr := fmt.Sprintf(":%d", port)
		logger.Log("transport", "HTTP", "addr", addr)
		errs <- http.ListenAndServe(addr, nil)
	}()

	logger.Log("exit", <-errs)
}
