package service

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	cb "github.com/callistaenterprise/goblog/common/circuitbreaker"

	"github.com/callistaenterprise/goblog/userservice/dbclient"
	"github.com/callistaenterprise/goblog/userservice/model"
	"github.com/gorilla/mux"
)

//DBClient db client
var DBClient dbclient.IBoltClient

var isHealthy = true

var client = &http.Client{}

var fallbackProduct = []model.Product{{
	Name:        "Default Product",
	Price:       0.0,
	Description: "Default description",
	ServedBy:    "circuit-breaker",
},
}

func init() {
	var transport http.RoundTripper = &http.Transport{
		DisableKeepAlives: true,
	}
	client.Transport = transport
}

//GetUser get user
func GetUser(w http.ResponseWriter, r *http.Request) {

	// Read the 'userID' path parameter from the mux map
	var userID = mux.Vars(r)["userID"]

	// Read the user struct BoltDB
	user, err := DBClient.QueryUser(userID)
	user.ServedBy = getIP()

	// If err, return a 404
	if err != nil {
		fmt.Println("Some error occured serving " + userID + ": " + err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// NEW call the productservice
	user.Products = getProducts()

	// If found, marshal into JSON, write headers and content
	data, _ := json.Marshal(user)
	writeJSONResponse(w, http.StatusOK, data)
}

func getProducts() []model.Product {

	body, err := cb.CallUsingCircuitBreaker("productservice", "http://productservice:8081/products/", "GET")
	if err == nil {
		products := []model.Product{}
		json.Unmarshal(body, &products)
		return products
	}

	return fallbackProduct
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Since we're here, we already know that HTTP service is up. Let's just check the state of the boltdb connection
	dbUp := DBClient.Check()
	if dbUp && isHealthy {
		data, _ := json.Marshal(healthCheckResponse{Status: "UP"})
		writeJSONResponse(w, http.StatusOK, data)
	} else {
		data, _ := json.Marshal(healthCheckResponse{Status: "Database unaccessible"})
		writeJSONResponse(w, http.StatusServiceUnavailable, data)
	}
}

func SetHealthyState(w http.ResponseWriter, r *http.Request) {
	// Read the 'state' path parameter from the mux map
	var state, err = strconv.ParseBool(mux.Vars(r)["state"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isHealthy = state
	w.WriteHeader(http.StatusOK)
}

func writeJSONResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

func getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	panic("Unable to determine local IP address (non loopback). Exiting.")
}

type healthCheckResponse struct {
	Status string `json:"status"`
}
