/*
 * Copyright (C) 2019 Lior Okman <lior.okman@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
		fmt.Printf("No VCAP_SERVICES")
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
