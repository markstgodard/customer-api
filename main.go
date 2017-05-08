package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
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

	// logging
	s = customers.NewLoggingMiddlware(logger)(s)

	// metrics
	fieldKeys := []string{"method"}
	s = customers.NewMetricsService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "customer_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "customer_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		s,
	)

	mux := http.NewServeMux()

	httpLogger := log.With(logger, "component", "http")

	mux.Handle("/api/v1/", customers.MakeHTTPHandler(s, httpLogger))
	http.Handle("/", mux)

	http.Handle("/metrics", promhttp.Handler())

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
