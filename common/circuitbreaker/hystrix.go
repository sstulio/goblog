package circuitbreaker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
)

//Client http client
var Client http.Client

//RETRIES Number of retries
var RETRIES = 3

//ConfigureHystrix config histrix
func ConfigureHystrix(commands []string) {

	for _, command := range commands {
		hystrix.ConfigureCommand(command, hystrix.CommandConfig{
			Timeout:                hystrix.DefaultTimeout,
			MaxConcurrentRequests:  hystrix.DefaultMaxConcurrent,
			ErrorPercentThreshold:  hystrix.DefaultErrorPercentThreshold,
			RequestVolumeThreshold: 3,
			SleepWindow:            hystrix.DefaultSleepWindow,
		})
	}
}

//CallUsingCircuitBreaker call using circuit breaker
func CallUsingCircuitBreaker(breakerName string, url string, method string) ([]byte, error) {
	output := make(chan []byte, 1)
	errors := hystrix.Go(breakerName, func() error {

		req, _ := http.NewRequest(method, url, nil)
		err := callWithRetries(req, output)

		return err
	}, func(err error) error {
		circuit, _, _ := hystrix.GetCircuit(breakerName)
		if circuit.IsOpen() {
			fmt.Printf("Circuit is Open!! \n")
		}
		return err
	})

	select {
	case out := <-output:
		fmt.Printf("Call in breaker %v successful \n", breakerName)
		return out, nil

	case err := <-errors:
		fmt.Printf("Got error on channel in breaker %v. Msg: %v \n", breakerName, err.Error())
		return nil, err
	}
}

func callWithRetries(req *http.Request, output chan []byte) error {

	r := retrier.New(retrier.ConstantBackoff(RETRIES, 100*time.Millisecond), nil)
	attempt := 0
	err := r.Run(func() error {
		attempt++
		resp, err := Client.Do(req)
		if err == nil && resp.StatusCode < 299 {
			responseBody, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				output <- responseBody
				return nil
			}
			return err
		} else if err == nil {
			err = fmt.Errorf("status was %v", resp.StatusCode)
		}

		return err
	})
	return err
}
