package service

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/callistaenterprise/goblog/productservice/model"
	"github.com/gorilla/mux"
)

var isHealthy = true

var client = &http.Client{}

func init() {
	var transport http.RoundTripper = &http.Transport{
		DisableKeepAlives: true,
	}
	client.Transport = transport
}

//GetProducts return products
func GetProducts(w http.ResponseWriter, r *http.Request) {

	if !isHealthy {
		data, _ := json.Marshal(healthCheckResponse{Status: "Unavailable"})
		writeJSONResponse(w, http.StatusServiceUnavailable, data)
		return
	}

	products := []model.Product{
		{
			Name:        "Smartphone Xiaomi",
			Description: "Super fast",
			Price:       1000.0,
			ServedBy:    getIP(),
		},
		{
			Name:        "SmartTV Samsumg",
			Description: "Comes with Netflix installed",
			Price:       2000.0,
			ServedBy:    getIP(),
		},
	}

	data, _ := json.Marshal(products)
	writeJSONResponse(w, http.StatusOK, data)
}

//HealthCheck health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if isHealthy {
		data, _ := json.Marshal(healthCheckResponse{Status: "UP"})
		writeJSONResponse(w, http.StatusOK, data)
	} else {
		data, _ := json.Marshal(healthCheckResponse{Status: "Unavailable"})
		writeJSONResponse(w, http.StatusServiceUnavailable, data)
	}
}

//SetHealthyState set healthy state
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
