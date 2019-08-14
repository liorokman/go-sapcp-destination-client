package main

import (
	"fmt"
	"os"

	destinations "github.com/liorokman/go-sapcp-destination-client"
	"github.com/tidwall/gjson"
)

func main() {

	vcap := os.Getenv("VCAP_SERVICES")
	if vcap == "" {
		fmt.Sprintf("No VCAP_SERVICES")
		os.Exit(1)
	}
	destinationClient, err := destinations.NewClient(destinations.DestinationClientConfiguration{
		ClientID:     gjson.Get(vcap, "destination.0.credentials.clientid").String(),
		ClientSecret: gjson.Get(vcap, "destination.0.credentials.clientsecret").String(),
		TokenURL:     gjson.Get(vcap, "destination.0.credentials.url").String(),
		ServiceURL:   gjson.Get(vcap, "destination.0.credentials.uri").String(),
	})
	if err != nil {
		panic(err)
	}
	destinations, err := destinationClient.GetSubaccountDestinations()
	if err != nil {
		panic(err)
	}
	for _, dest := range destinations {
		fmt.Printf("Destination: %s(%s)\n", dest.Name, dest.Type)
	}

}
