package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var (
	addr = flag.String("listen-address", ":8080", "The address for the http server to listen on")

	currentFlag float64 = 0

	FlagMetric = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "service_pending_queries",
			Help: "User controlled flag.  Settable at /flag",
		},
	)

	SimpleCounter = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "service",
			Name:      "requested_in_porgress",
			Help:      "User controlled flag.  Settable at /flag",
		},
	)
)

func init() {
	prometheus.MustRegister(FlagMetric)
	prometheus.MustRegister(SimpleCounter)
	//prometheus.MustRegister(SimpleGague)
	//prometheus.MustRegister(SimpleSummary)
}

func updateFlag(w http.ResponseWriter, r *http.Request) {
	var n struct {
		Number float64 `json:"number"`
	}

	if r.Method == "GET" {
		fmt.Fprintf(w, "{\"number\": %f}", currentFlag)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	if r.Body == nil {
		http.Error(w, "Please include a number", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("internal error", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &n)
	if err != nil {
		log.Error("incorrect json", "error", err)
		http.Error(w, "incorrect json", http.StatusBadRequest)
		return
	}

	currentFlag = n.Number
	FlagMetric.Set(currentFlag)

	fmt.Fprintf(w, "{\"success\": true,\"number\": %f}", currentFlag)
}

func main() {

	flag.Parse()

	FlagMetric.Set(currentFlag)

	http.HandleFunc("/flag", updateFlag)
	http.Handle("/metrics", promhttp.Handler())
	log.Infof("Beging http server", "listen-address", *addr)
	log.Fatalf("http serv error", "error", http.ListenAndServe(*addr, nil))
}
