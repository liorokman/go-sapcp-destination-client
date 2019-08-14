package gosapcpdestinationclient

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2/clientcredentials"
)

type destinationClient struct {
	restyClient *resty.Client
}

type DestinationClient interface {
	DestinationFinder

	SubaccountDestinationManager
	SubaccountCertificateManager

	InstanceDestinationManager
	InstanceCertificateManager

	SetDebug(debug bool)
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

func NewClient(clientConf DestinationClientConfiguration) (DestinationClient, error) {
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

	return &destinationClient{
		restyClient: restyClient,
	}, nil
}

/****************************   Find a destination **********************************/

func (d *destinationClient) Find(name string) (DestinationLookupResult, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

/**************************** Destinatons on a subaccount level **********************************/

func (d *destinationClient) GetSubaccountDestinations() ([]Destination, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

func (d *destinationClient) CreateSubaccountDestination(newDestination Destination) error {

	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetBody(newDestination).
		SetError(&errResponse).
		Post("/subaccountDestinations")

	if err != nil {
		return err
	}
	if response.StatusCode() != 201 {
		return errResponse
	}
	return nil
}

func (d *destinationClient) UpdateSubaccountDestination(dest Destination) (AffectedRecords, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

func (d *destinationClient) GetSubaccountDestination(name string) (Destination, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

func (d *destinationClient) DeleteSubaccountDestination(name string) (AffectedRecords, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

/**************************** Subaccount Certificates **********************************/

func (d *destinationClient) GetSubaccountCertificates() ([]Certificate, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

func (d *destinationClient) CreateSubaccountCertificate(cert Certificate) error {

	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetBody(cert).
		SetError(&errResponse).
		Post("/subaccountCertificates")

	if err != nil {
		return err
	}
	if response.StatusCode() != 201 {
		return errResponse
	}
	return nil
}

func (d *destinationClient) GetSubaccountCertificate(name string) (Certificate, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

func (d *destinationClient) DeleteSubaccountCertificate(name string) (AffectedRecords, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

/**************************** Destinatons on an instance level **********************************/

func (d *destinationClient) GetInstanceDestinations() ([]Destination, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

func (d *destinationClient) CreateInstanceDestination(newDestination Destination) error {

	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetBody(newDestination).
		SetError(&errResponse).
		Post("/instanceDestinations")

	if err != nil {
		return err
	}
	if response.StatusCode() != 201 {
		return errResponse
	}
	return nil
}

func (d *destinationClient) UpdateInstanceDestination(dest Destination) (AffectedRecords, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

func (d *destinationClient) GetInstanceDestination(name string) (Destination, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

func (d *destinationClient) DeleteInstanceDestination(name string) (AffectedRecords, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

/**************************** Instance Certificates **********************************/

func (d *destinationClient) GetInstanceCertificates() ([]Certificate, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

func (d *destinationClient) CreateInstanceCertificate(cert Certificate) error {

	var errResponse ErrorMessage

	response, err := d.restyClient.R().
		SetBody(cert).
		SetError(&errResponse).
		Post("/instanceCertificates")

	if err != nil {
		return err
	}
	if response.StatusCode() != 201 {
		return errResponse
	}
	return nil
}

func (d *destinationClient) GetInstanceCertificate(name string) (Certificate, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

func (d *destinationClient) DeleteInstanceCertificate(name string) (AffectedRecords, error) {

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
		return retval, errResponse
	}
	return retval, nil
}

/****************************** Misc. ************************************************/

func (d *destinationClient) SetDebug(debug bool) {
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

func (e ErrorMessage) Error() string {
	return e.ErrorMessage
}
