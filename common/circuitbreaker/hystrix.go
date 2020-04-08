package circuitbreaker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
)

var Client http.Client

var RETRIES = 3

func ConfigureHystrix(commands []string) {

	for _, command := range commands {
		hystrix.ConfigureCommand(command, hystrix.CommandConfig{
			Timeout:                hystrix.DefaultTimeout,
			MaxConcurrentRequests:  hystrix.DefaultMaxConcurrent,
			ErrorPercentThreshold:  hystrix.DefaultErrorPercentThreshold,
			RequestVolumeThreshold: hystrix.DefaultVolumeThreshold,
			SleepWindow:            hystrix.DefaultSleepWindow,
		})
	}
}

func CallUsingCircuitBreaker(breakerName string, url string, method string) ([]byte, error) {
	output := make(chan []byte, 1)
	errors := hystrix.Go(breakerName, func() error {

		req, _ := http.NewRequest(method, url, nil)
		err := callWithRetries(req, output)

		return err // For hystrix, forward the err from the retrier. It's nil if OK.
	}, func(err error) error {
		fmt.Printf("In fallback function for breaker %v, error: %v \n", breakerName, err.Error())
		circuit, _, _ := hystrix.GetCircuit(breakerName)
		fmt.Printf("Circuit state is: %v \n", circuit.IsOpen())
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
			err = fmt.Errorf("Status was %v \n", resp.StatusCode)
		}

		fmt.Printf("Retrier failed, attempt %v \n", attempt)

		return err
	})
	return err
}
