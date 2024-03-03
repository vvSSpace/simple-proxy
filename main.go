package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"simple-proxy/methods/api"
	"simple-proxy/methods/proxy"
	"strings"
)

var proxyPort int
var omUrl string
var lcUrl string
var newCustomerID int64
var oldCustomerID int64
var CustomerIDForHistory int64

func init() {
	flag.IntVar(&proxyPort, "p", 8885, "Proxy port number")
	flag.StringVar(&omUrl, "om", "http://orders-test.stage2.com", "Orders service URL")
	flag.StringVar(&lcUrl, "lc", "http://lc-test.stage2.com", "Lc service URL")
	flag.Int64Var(&newCustomerID, "c", 0, "Customer ID to replace")
	flag.Int64Var(&CustomerIDForHistory, "ch", 0, "Customer ID to request lc history")
	flag.Parse()

	if proxyPort == 0 || omUrl == "" || lcUrl == "" || newCustomerID == 0 || CustomerIDForHistory == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	condition := make(chan bool)

	http.HandleFunc("/v2/orders.get-order", func(w http.ResponseWriter, r *http.Request) {
		proxy.OrdersGetOrder(w, r, condition, omUrl, &oldCustomerID, newCustomerID)
	})

	http.HandleFunc("/v1/orders/get-info", func(w http.ResponseWriter, r *http.Request) {
		proxy.LcGetInfo(w, r, condition, lcUrl, &oldCustomerID)
	})

	http.HandleFunc("/v1/points/get-history", func(w http.ResponseWriter, r *http.Request) {
		proxy.LcGetHistory(w, r, lcUrl, CustomerIDForHistory)
	})

	http.HandleFunc("/update_customer_id", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			log.Println("Received request: POST", r.URL)
			api.UpdateCustomerForHistory(w, r, &CustomerIDForHistory)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/__settings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			log.Println("Received request: GET", r.URL)
			api.GetSettings(w, omUrl, lcUrl, newCustomerID, CustomerIDForHistory)
		case http.MethodPost:
			log.Println("Received request: POST", r.URL)
			api.PostSettings(w, r, &omUrl, &lcUrl, &newCustomerID, &CustomerIDForHistory)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	address := fmt.Sprintf(":%d", proxyPort)
	host := "http://localhost" + address
	log.Println(strings.Repeat("-", 80))
	log.Println("Starting proxy on:                ", host)
	log.Println("Orders service URL:               ", omUrl)
	log.Println("Lc service URL:                   ", lcUrl)
	log.Println("Customer ID to replace (orders):  ", newCustomerID)
	log.Println("Customer ID to replace (history): ", CustomerIDForHistory)
	log.Println(strings.Repeat("-", 80))
	log.Fatal(http.ListenAndServe(address, nil))
}
