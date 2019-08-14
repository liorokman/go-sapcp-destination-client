# Golang client for SAP Cloud Platform Destination Services

[![GoDoc](https://godoc.org/github.com/liorokman/go-sapcp-destination-client?status.svg)](https://godoc.org/github.com/liorokman/go-sapcp-destination-client) 
[![Go Report Card](https://goreportcard.com/badge/github.com/liorokman/go-sapcp-destination-client)](https://goreportcard.com/report/github.com/liorokman/go-sapcp-destination-client)

Based on the published API at https://api.sap.com/api/SAP_CP_CF_Connectivity_Destination/resource.


This library provides a convenient client for accessing the Destination service on the SAP Cloud Platform Cloud Foundry environments.

## Usage

1. The configuration details required for creating a new `DestinationClient` instance are provided in the destination service binding in the `credentials` section. The application must 
create an instance of the destination service, and bind to it in order to access the configuration details.

   ```bash
   cf create-service destination lite example-destination
   ```

1. Add the service to the `manifest.yml` in order to bind to it.
1. Configure and use the `DestinationClient`

```golang
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
```




