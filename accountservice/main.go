package main

import (
	"fmt"

	cb "github.com/callistaenterprise/goblog/common/circuitbreaker"

	"github.com/callistaenterprise/goblog/accountservice/dbclient"
	"github.com/callistaenterprise/goblog/accountservice/service"
)

var appName = "accountservice"

func main() {
	fmt.Printf("Starting %v\n", appName)

	initializeBoltClient()
	cb.ConfigureHystrix([]string{"quoteservice"})

	service.StartWebServer("8080")
}

func initializeBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}
