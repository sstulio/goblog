package main

import (
	"fmt"

	"github.com/callistaenterprise/goblog/productservice/service"
)

var appName = "productservice"

func main() {
	fmt.Printf("Starting %v\n", appName)
	service.StartWebServer("8081")
}
