package gosapcpdestinationclient

import (
	"fmt"
)

func ExampleNewClient() {

	client, err := NewClient(DestinationClientConfiguration{
		ClientID:     "clientid",
		ClientSecret: "clientsecret",
		TokenURL:     "https://subdomain.authentication.eu10.hana.ondemand.com",
		ServiceURL:   "https://destination-configuration.cfapps.eu10.hana.ondemand.com",
	})
	if err != nil {
		panic(err)
	}
	destinations, err := client.GetSubaccountDestinations()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", destinations)
}
