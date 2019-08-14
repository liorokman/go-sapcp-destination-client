/*
Copyright (C) 2019 Lior Okman <lior.okman@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
