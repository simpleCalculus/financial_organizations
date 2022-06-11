package main

import (
	"financial_organizations/app"
	"financial_organizations/config"
	"log"
)

func main() {
	appInstance, err := app.NewApp(config.New())
	if err != nil {
		log.Fatalf("fauled to create a new app %v", err)
	}
	appInstance.Run(":8080")
}
