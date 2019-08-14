package gosapcpdestinationclient

// Types used by the RESTful API

// DestinationType enumeration
type DestinationType string

const (
	// HTTP destination
	HTTPDestination DestinationType = "HTTP"
	// RFC destination
	RFCDestination DestinationType = "RFC"
	// Mail (SMTP) destination
	MailDestination DestinationType = "MAIL"
	// LDAP destination
	LDAPDestination DestinationType = "LDAP"
)

// ErrorMessage struct contains errors returned by the Destination API
type ErrorMessage struct {
	ErrorMessage string `json:"ErrorMessage"`
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
