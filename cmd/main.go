package main

import (
	"flag"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rguilmont/cypress-dashboard-exporter/pkg/cypresscollector"
	"github.com/sirupsen/logrus"
)

const readinessTime = 360

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logrus.Info("%v ")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func initCollector(u *url.URL, project, email, password string, keepUntil int64) *cypresscollector.CypressDashboardCollector {
	cypressCollector, err := cypresscollector.NewCypressDashboardCollector(*u, project, email, password, keepUntil)
	if err != nil {
		logrus.Panicln(err)
	}
	return cypressCollector
}

// Convert the number of days into seconds
func toSeconds(days int64) int64 {
	return days * int64(time.Hour) * 24
}

func main() {
	listen := flag.String("listen", "0.0.0.0:8081", "host:port to listen")
	project := flag.String("project", "7s5okt", "host:port to listen")
	keepUntil := flag.Int64("keepUntil", 14,
		"Time ( in days ) to keep in memory the results of a test/run before removing it.")

	email := flag.String("email", "", "email to connect to the dashboard")
	password := flag.String("password", "", "password to connect to the dashboard")
	debug := flag.Bool("debug", false, "activate debug logging")

	flag.Parse()
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Info("Starting Cypress dashboard exporter")
	parsedURL, err := url.Parse("https://dashboard.cypress.io/graphql")

	if err != nil {
		logrus.Panicln("Impossible to parse URL ", err)
	}

	ddCollector := initCollector(parsedURL, *project, *email, *password, toSeconds(*keepUntil))
	logrus.Infoln("Monitoring Cypress dashboard at ", parsedURL, "for project ID ", *project)
	logrus.Infoln("Keeping old timeseries for %v days ", *keepUntil)
	prometheus.MustRegister(ddCollector)
	http.Handle("/metrics", promhttp.Handler())

	logrus.Info("Listening ", *listen)
	logrus.Fatal(http.ListenAndServe(*listen, handlers.LoggingHandler(logrus.New().Out, http.DefaultServeMux)))
}
