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
	"context"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2/clientcredentials"
)

// DestinationClient provides the client object for accessing destinations in the SAP Cloud Platform Cloud Foundry environment.
type DestinationClient struct {
	restyClient *resty.Client
}

// DestinationFinder provides a Find method for discovering destinations on any level.
type DestinationFinder interface {
	Find(name string) (DestinationLookupResult, error)
}

// SubaccountDestinationManager provides an interface for methods that manage destinations on the Subaccount level
type SubaccountDestinationManager interface {
	GetSubaccountDestinations() ([]Destination, error)
	CreateSubaccountDestination(newDestination Destination) error
	UpdateSubaccountDestination(dest Destination) (AffectedRecords, error)
	GetSubaccountDestination(name string) (Destination, error)
	DeleteSubaccountDestination(name string) (AffectedRecords, error)
}

// SubaccountCertificateManager provides an interface for methods that manage certificates on the Subaccount level
type SubaccountCertificateManager interface {
	GetSubaccountCertificates() ([]Certificate, error)
	CreateSubaccountCertificate(cert Certificate) error
	GetSubaccountCertificate(name string) (Certificate, error)
	DeleteSubaccountCertificate(name string) (AffectedRecords, error)
}

// InstanceDestinationManager provides an interface for methods that manage destinations on the Instance level
type InstanceDestinationManager interface {
	GetInstanceDestinations() ([]Destination, error)
	CreateInstanceDestination(newDestination Destination) error
	UpdateInstanceDestination(dest Destination) (AffectedRecords, error)
	GetInstanceDestination(name string) (Destination, error)
	DeleteInstanceDestination(name string) (AffectedRecords, error)
}

// InstanceCertificateManager provides an interface for methods that manage certificates on the Instance level
type InstanceCertificateManager interface {
	GetInstanceCertificates() ([]Certificate, error)
	CreateInstanceCertificate(cert Certificate) error
	GetInstanceCertificate(name string) (Certificate, error)
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

// NewClient creates a new DestinationClient object configured according to the provided DestinationClientConfiguration object
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

// Find a destination by name on all levels and return the first match.
// Search priority is destination on service instance level. If none is found, fallbacks to subaccount level (accessible by all apps deployed in the same subaccount).
// If userToken is not empty, it is passed as the value of the `X-user-token` header. This enables token-exchange flows via the Find operation. If a token-exchange
// is not required, pass an empty string as the userToken value.
func (d *DestinationClient) Find(name string, userToken string) (DestinationLookupResult, error) {

	var retval DestinationLookupResult
	var errResponse ErrorMessage

	request := d.restyClient.R().
		SetResult(&retval).
		SetError(&errResponse).
		SetPathParams(map[string]string{
			"name": name,
		})
	if userToken != "" {
		request.SetHeader("X-user-token", userToken)
	}
	response, err := request.Get("/destinations/{name}")

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

// GetSubaccountDestinations returns a list of destinations posted on subaccount level. If none is found, an empty array is returned. Subaccount is determined by the passed OAuth access token.
func (d *DestinationClient) GetSubaccountDestinations() ([]Destination, error) {

	var retval = make([]Destination, 0)
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

// CreateSubaccountDestination creates a new destination on subaccount level. Subaccount is determined by the passed OAuth access token.
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

// UpdateSubaccountDestination updates (overwrites) an existing destination with a new destination, posted on subaccount level. Subaccount is determined by the passed OAuth access token
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

// GetSubaccountDestination retrieves a named destination posted on subaccount level. Subaccount is determined by the passed OAuth access token.
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

// DeleteSubaccountDestination deletes a destination posted on subaccount level. Subaccount is determined by the passed OAuth access token.
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

// GetSubaccountCertificates retrieves all certificates posted on the subaccount level. In none are found, an empty array is returned. The Subaccount is determined based on the passed OAuth access token
func (d *DestinationClient) GetSubaccountCertificates() ([]Certificate, error) {

	var retval = make([]Certificate, 0)
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

// CreateSubaccountCertificate creates a new certificate on the subaccount level. The Subaccount is determined by the passed OAuth access token
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

// GetSubaccountCertificate retrieves a named certificate posted on the subaccount level. The Subaccount is determined by the passed OAuth access token
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

// DeleteSubaccountCertificate deletes a certificate posted on the subaccount level. The Subaccount is determined by the passed OAuth access token
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

// GetInstanceDestinations retrieves all destinations on the service instance level. If none are found, an empty list is returned. Service instance and subaccount are determined the passed OAuth access token
func (d *DestinationClient) GetInstanceDestinations() ([]Destination, error) {

	var retval = make([]Destination, 0)
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

// CreateInstanceDestination creates a new destination on the service instance level. The service instance and subaccount are determined by the passed OAuth access token
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

// UpdateInstanceDestination updates (overwrites) an existing destination with the passed destination. The service instance and subaccount are determined by the passed OAuth access token
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

// GetInstanceDestination retrieves a destination posted on the service instance level. The service instance and subaccount are determined by the passed OAuth access token
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

// DeleteInstanceDestination deletes a destination posted on the service instance level. The service instance and subaccount are determined by the passed OAuth access token
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

// GetInstanceCertificates retrieves all certificates posted on the service instance level. If none are found, an empty list is returned. The service instance and subaccount are determined by the passed OAuth access token
func (d *DestinationClient) GetInstanceCertificates() ([]Certificate, error) {

	var retval = make([]Certificate, 0)
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

// CreateInstanceCertificate creates a new certificate on the service instance level. The service instance and subaccount are determined by the passed OAuth access token
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

// GetInstanceCertificate retrieves a certificate posted on the service instance level. The service instance and subaccount are determined by the passed OAuth access token
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

// DeleteInstanceCertificate deletes a certificate posted on the service instance level. The service instance and subaccount are determined by the passed OAuth access token
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

// SetDebug enables or disables debug output for the DestinationClient
func (d *DestinationClient) SetDebug(debug bool) {
	d.restyClient.SetDebug(debug)
}

// MarshalJSON marshalls a Destination object as expected by the Destination RESTful API
func (d Destination) MarshalJSON() ([]byte, error) {
	if d.Properties == nil {
		d.Properties = make(map[string]string)
	}
	d.Properties["Name"] = d.Name
	d.Properties["Type"] = string(d.Type)
	return json.Marshal(d.Properties)
}

// UnmarshalJSON unmarshalls a Destination object as provided by the Destination RESTful API
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

// StatusCode returns the status code provided with the error
func (e ErrorMessage) StatusCode() int {
	return e.statusCode
}

func (e ErrorMessage) Error() string {
	return e.ErrorMessage
}
