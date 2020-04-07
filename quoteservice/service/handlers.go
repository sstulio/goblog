package service

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/callistaenterprise/goblog/quoteservice/model"
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

func GetQuote(w http.ResponseWriter, r *http.Request) {

	quote := model.Quote{
		Text:     "Be or not to be",
		Language: "English",
	}
	quote.ServedBy = getIP()

	// If found, marshal into JSON, write headers and content
	data, _ := json.Marshal(quote)
	writeJsonResponse(w, http.StatusOK, data)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if isHealthy {
		data, _ := json.Marshal(healthCheckResponse{Status: "UP"})
		writeJsonResponse(w, http.StatusOK, data)
	} else {
		data, _ := json.Marshal(healthCheckResponse{Status: "Unavailable"})
		writeJsonResponse(w, http.StatusServiceUnavailable, data)
	}
}

func SetHealthyState(w http.ResponseWriter, r *http.Request) {
	// Read the 'accountId' path parameter from the mux map
	var state, err = strconv.ParseBool(mux.Vars(r)["state"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isHealthy = state
	w.WriteHeader(http.StatusOK)
}

func writeJsonResponse(w http.ResponseWriter, status int, data []byte) {
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
