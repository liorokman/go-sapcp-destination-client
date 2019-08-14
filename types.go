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

// Types used by the RESTful API

// DestinationType enumeration
type DestinationType string

const (
	// HTTPDestination type
	HTTPDestination DestinationType = "HTTP"
	// RFCDestination type
	RFCDestination DestinationType = "RFC"
	// MailDestination type
	MailDestination DestinationType = "MAIL"
	// LDAPDestination type
	LDAPDestination DestinationType = "LDAP"
)

// Common destination properties
const (
	// Property name for destination Description
	DescriptionProperty = "Description"
	// Property name for the destination Authentication property
	AuthenticationProperty = "Authentication"

	// Valid values for the authentication property

	AppToAppSSOAuthentication               = "AppToAppSSO"
	BasicAuthentication                     = "BasicAuthentication"
	ClientCertificateAuthentication         = "ClientCertificateAuthentication"
	NoAuthentication                        = "NoAuthentication"
	OAuth2ClientCredentialsAuthentication   = "OAuth2ClientCredentials"
	OAuth2SAMLBearerAssertionAuthentication = "OAuth2SAMLBearerAssertion"
	OAuth2UserTokenExchangeAuthentication   = "OAuth2UserTokenExchange"
	SAPAssetionSSOAuthentication            = "SAPAssertionSSO"

	// Property name for the destination ProxyType property
	ProxyTypeProperty = "ProxyType"

	// Valid values for the ProxyType property
	InternetProxy  = "Internet"
	OnPremiseProxy = "OnPremise"

	// Property name for the destination URL property
	URLProperty = "URL"

	// Property name for the LocationID destination property
	LocationIDProperty = "LocationID"

	// Property name for the destination User property
	UserProperty = "User"

	// Property name for the destination Password property
	PasswordProperty = "Password"

	// Property name for the destination RepositoryUser property
	RepoUserProperty = "RepositoryUser"

	// Property name for the destination RepositoryPassword property
	RepoPasswordProperty = "RepositoryPassword"
)

// ErrorMessage struct contains errors returned by the Destination API
type ErrorMessage struct {
	ErrorMessage string `json:"ErrorMessage"`
	statusCode   int
}

// Destination describes a single Destination
type Destination struct {
	// The name of the destination
	Name string
	// The type of the destination
	Type DestinationType
	// Any properties defined on the destination
	Properties map[string]string
}

// Certificate describes a single certificate
type Certificate struct {
	// The name of the destination
	Name string `json:"Name"`
	// The type of the destination
	Type string `json:"Type"`
	// Base64 encoded keystore/certificate binary content
	Content string `json:"Content"`
}

// AuthToken describes an authentication token
type AuthToken struct {
	// Type of the authentication token
	Type string `json:"type"`
	// Value of the authentication token
	Value string `json:"value"`
}

// Owner describes the level on which the destination is defined.
// At least one of SubaccountID or InstanceID are guaranteed to have a value.
type Owner struct {
	// Subaccount ID owning this destination
	SubaccountID string `json:"SubaccountId,omitempty"`
	// Instance ID owning this destination
	InstanceID string `json:"InstanceId,omitempty"`
}

// DestinationLookupResult contains the result of a find operation
type DestinationLookupResult struct {
	// The level on which the destination is defined
	Owner Owner `json:"owner,omitempty"`
	// The destination information
	Destination Destination `json:"destinationConfiguration,omitempty"`
	// Certificates (if present) for the destination
	Certificates []Certificate `json:"certificates,omitempty"`
	// Authentication tokens (if present) for the destination
	AuthTokens []AuthToken `json:"authTokens,omitempty"`
}

// AffectedRecords contains the number of records affected by the operation
type AffectedRecords struct {
	Count int `json:"count"`
}
