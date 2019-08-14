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
package gosapcpdestinationclient

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2/clientcredentials"
)

type DestinationClient struct {
	restyClient *resty.Client
}

type DestinationFinder interface {
	// Find a destination by name on all levels and return the first match.
	// Search priority is destination on service instance level. If none is found, fallbacks to subaccount level (accessible by all apps deployed in the same subaccount).
	Find(name string) (DestinationLookupResult, error)
}

type SubaccountDestinationManager interface {
	// Get a list of destinations posted on subaccount level. If none is found, an empty array is returned. Subaccount is determined by the passed OAuth access token.
	GetSubaccountDestinations() ([]Destination, error)

	// Create a new destination on subaccount level. Subaccount is determined by the passed OAuth access token.
	CreateSubaccountDestination(newDestination Destination) error

	// Update (overwrite) existing destination with a new destination, posted on subaccount level. Subaccount is determined by the passed OAuth access token
	UpdateSubaccountDestination(dest Destination) (AffectedRecords, error)

	// Get a destination posted on subaccount level. Subaccount is determined by the passed OAuth access token.
	GetSubaccountDestination(name string) (Destination, error)

	// Delete a destination posted on subaccount level. Subaccount is determined by the passed OAuth access token.
	DeleteSubaccountDestination(name string) (AffectedRecords, error)
}

type SubaccountCertificateManager interface {
	// Get all certificates posted on subaccount level. In none is found, an empty array is returned. Subaccount is determined by the passed OAuth access token
	GetSubaccountCertificates() ([]Certificate, error)

	// Create a new certificate on subaccount level. Subaccount is determined by the passed OAuth access token
	CreateSubaccountCertificate(cert Certificate) error

	// Get a certificate posted on subaccount level. Subaccount is determined by the passed OAuth access token
	GetSubaccountCertificate(name string) (Certificate, error)

	// Delete a certificate posted on subaccount level. Subaccount is determined by the passed OAuth access token
	DeleteSubaccountCertificate(name string) (AffectedRecords, error)
}

type InstanceDestinationManager interface {
	// Get all destinations on service instance level. If none is found, an empty list is returned. Service instance and subaccount are determined the passed OAuth access token
	GetInstanceDestinations() ([]Destination, error)

	// Create a new destination on service instance level. Service instance and subaccount are determined by the passed OAuth access token
	CreateInstanceDestination(newDestination Destination) error

	// Update (overwrite) the existing destination with the passed destination. Service instance and subaccount are determined by the passed OAuth access token
	UpdateInstanceDestination(dest Destination) (AffectedRecords, error)

	// Get a destination posted on service instance level. Service instance and subaccount are determined by the passed OAuth access token
	GetInstanceDestination(name string) (Destination, error)

	// Delete a destination posted on service instance level. Service instance and subaccount are determined by the passed OAuth access token
	DeleteInstanceDestination(name string) (AffectedRecords, error)
}

type InstanceCertificateManager interface {
	// Get all certificates posted on service instance level. If none is found, an empty list is returned. Service instance and subaccount are determined by the passed OAuth access token
	GetInstanceCertificates() ([]Certificate, error)

	// Create a new certificate on service instance level. Service instance and subaccount are determined by the passed OAuth access token
	CreateInstanceCertificate(cert Certificate) error

	// Get a certificate posted on service instance level. Service instance and subaccount are determined by the passed OAuth access token
	GetInstanceCertificate(name string) (Certificate, error)

	// Deletes a certificate posted on service instance level. Service instance and subaccount are determined by the passed OAuth access token
	DeleteInstanceCertificate(name string) (AffectedRecords, error)
}

// DestinationClientConfiguration contains the values required for configuring a new Destination client
type DestinationClientConfiguration struct {
	// ClientID for authentication purposes. Use the clientid attribute in the service binding
	ClientID string
	// ClientSecret for authentication purposes. Use the clientsecret attribute in the service binding
	ClientSecret string
	// TokenURL for authentication purposes. Use the url attribute in the service binding
	TokenURL string
	// ServiceURL for accessing the service RESTful endpoint. Use the uri attribute in the service binding
	ServiceURL string
}

func NewClient(clientConf DestinationClientConfiguration) (*DestinationClient, error) {
	conf := &clientcredentials.Config{
		ClientID:     clientConf.ClientID,
		ClientSecret: clientConf.ClientSecret,
		TokenURL:     clientConf.TokenURL + "/oauth/token",
		Scopes:       []string{},
	}
	client := conf.Client(context.Background())

	restyClient := resty.NewWithClient(client).
		SetHostURL(clientConf.ServiceURL+"/destination-configuration/v1").
		SetHeader("Accept", "application/json").
		SetTimeout(60 * time.Second)

	return &DestinationClient{
		restyClient: restyClient,
	}, nil
}

/****************************   Find a destination **********************************/

func (d *DestinationClient) Find(name string) (DestinationLookupResult, error) {

	var retval DestinationLookupResult
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		SetPathParams(map[string]string{
			"name": name,
		}).
		Get("/destinations/{name}")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

/**************************** Destinations on a subaccount level **********************************/

func (d *DestinationClient) GetSubaccountDestinations() ([]Destination, error) {

	var retval []Destination = make([]Destination, 0)
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		Get("/subaccountDestinations")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

func (d *DestinationClient) CreateSubaccountDestination(newDestination Destination) error {

	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetBody(newDestination).
		SetError(&errResponse).
		Post("/subaccountDestinations")

	if err != nil {
		return err
	}
	if response.StatusCode() != 201 {
		errResponse.statusCode = response.StatusCode()
		return errResponse
	}
	return nil
}

func (d *DestinationClient) UpdateSubaccountDestination(dest Destination) (AffectedRecords, error) {

	var retval AffectedRecords
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetBody(dest).
		SetResult(&retval).
		SetError(&errResponse).
		Put("/subaccountDestinations")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

func (d *DestinationClient) GetSubaccountDestination(name string) (Destination, error) {

	var retval Destination
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		SetPathParams(map[string]string{
			"name": name,
		}).
		Get("/subaccountDestinations/{name}")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

func (d *DestinationClient) DeleteSubaccountDestination(name string) (AffectedRecords, error) {

	var retval AffectedRecords
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		SetPathParams(map[string]string{
			"name": name,
		}).
		Delete("/subaccountDestinations/{name}")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

/**************************** Subaccount Certificates **********************************/

func (d *DestinationClient) GetSubaccountCertificates() ([]Certificate, error) {

	var retval []Certificate = make([]Certificate, 0)
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		Get("/subaccountCertificates")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

func (d *DestinationClient) CreateSubaccountCertificate(cert Certificate) error {

	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetBody(cert).
		SetError(&errResponse).
		Post("/subaccountCertificates")

	if err != nil {
		return err
	}
	if response.StatusCode() != 201 {
		errResponse.statusCode = response.StatusCode()
		return errResponse
	}
	return nil
}

func (d *DestinationClient) GetSubaccountCertificate(name string) (Certificate, error) {

	var retval Certificate
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		SetPathParams(map[string]string{
			"name": name,
		}).
		Get("/subaccountCertificate/{name}")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

func (d *DestinationClient) DeleteSubaccountCertificate(name string) (AffectedRecords, error) {

	var retval AffectedRecords
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		SetPathParams(map[string]string{
			"name": name,
		}).
		Delete("/subaccountCertificate/{name}")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

/**************************** Destinations on an instance level **********************************/

func (d *DestinationClient) GetInstanceDestinations() ([]Destination, error) {

	var retval []Destination = make([]Destination, 0)
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		Get("/instanceDestinations")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

func (d *DestinationClient) CreateInstanceDestination(newDestination Destination) error {

	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetBody(newDestination).
		SetError(&errResponse).
		Post("/instanceDestinations")

	if err != nil {
		return err
	}
	if response.StatusCode() != 201 {
		errResponse.statusCode = response.StatusCode()
		return errResponse
	}
	return nil
}

func (d *DestinationClient) UpdateInstanceDestination(dest Destination) (AffectedRecords, error) {

	var retval AffectedRecords
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetBody(dest).
		SetResult(&retval).
		SetError(&errResponse).
		Put("/instanceDestinations")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

func (d *DestinationClient) GetInstanceDestination(name string) (Destination, error) {

	var retval Destination
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		SetPathParams(map[string]string{
			"name": name,
		}).
		Get("/instanceDestinations/{name}")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

func (d *DestinationClient) DeleteInstanceDestination(name string) (AffectedRecords, error) {

	var retval AffectedRecords
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		SetPathParams(map[string]string{
			"name": name,
		}).
		Delete("/instanceDestinations/{name}")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

/**************************** Instance Certificates **********************************/

func (d *DestinationClient) GetInstanceCertificates() ([]Certificate, error) {

	var retval []Certificate = make([]Certificate, 0)
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		Get("/instanceCertificates")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

func (d *DestinationClient) CreateInstanceCertificate(cert Certificate) error {

	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetBody(cert).
		SetError(&errResponse).
		Post("/instanceCertificates")

	if err != nil {
		return err
	}
	if response.StatusCode() != 201 {
		errResponse.statusCode = response.StatusCode()
		return errResponse
	}
	return nil
}

func (d *DestinationClient) GetInstanceCertificate(name string) (Certificate, error) {

	var retval Certificate
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		SetPathParams(map[string]string{
			"name": name,
		}).
		Get("/instanceCertificate/{name}")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

func (d *DestinationClient) DeleteInstanceCertificate(name string) (AffectedRecords, error) {

	var retval AffectedRecords
	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		SetPathParams(map[string]string{
			"name": name,
		}).
		Delete("/instanceCertificate/{name}")

	if err != nil {
		return retval, err
	}
	if response.StatusCode() != 200 {
		errResponse.statusCode = response.StatusCode()
		return retval, errResponse
	}
	return retval, nil
}

/****************************** Misc. ************************************************/

func (d *DestinationClient) SetDebug(debug bool) {
	d.restyClient.SetDebug(debug)
}

func (d Destination) MarshalJSON() ([]byte, error) {
	d.Properties["Name"] = d.Name
	d.Properties["Type"] = string(d.Type)
	return json.Marshal(d.Properties)
}

func (d *Destination) UnmarshalJSON(b []byte) error {

	unmarshalled := map[string]string{}
	if err := json.Unmarshal(b, &unmarshalled); err != nil {
		return err
	}
	d.Properties = make(map[string]string)
	for k, v := range unmarshalled {
		switch k {
		case "Name":
			d.Name = v
		case "Type":
			switch v {
			case "HTTP":
				d.Type = HTTPDestination
			case "RFC":
				d.Type = RFCDestination
			case "MAIL":
				d.Type = MailDestination
			case "LDAP":
				d.Type = LDAPDestination
			default:
				d.Type = ""
			}
		default:
			d.Properties[k] = v
		}
	}
	return nil
}

func (e ErrorMessage) StatusCode() int {
	return e.statusCode
}

func (e ErrorMessage) Error() string {
	return e.ErrorMessage
}
