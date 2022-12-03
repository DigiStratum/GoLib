package identityprovider

/*
Identity Provider Service layer - here we set out the expected interface for identity provider
implementations. An IdentityProvider may be leverages in as simple or as intricate a way as one can
imagine ranging from Basic auth implementation to Oauth2.

The basic interface provides for a unique IDP ID that differentiates each from one another in a
collection of IDP's and other refernces. There is also a provision for a username+password check.
*/

type IdentityProviderIfc interface {
	GetId() string
	CheckCredentials(username, password string) bool
}
