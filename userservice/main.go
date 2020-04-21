package main

import (
	"fmt"

	cb "github.com/callistaenterprise/goblog/common/circuitbreaker"

	"github.com/callistaenterprise/goblog/userservice/dbclient"
	"github.com/callistaenterprise/goblog/userservice/service"
)

var appName = "userservice"

func main() {
	fmt.Printf("Starting %v\n", appName)

	initializeBoltClient()
	cb.ConfigureHystrix([]string{"productservice"})

	service.StartWebServer("8080")
}

func initializeBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}
