package asymmetric

/*

This is an opinionated implementation of asymmetric cryptography to keep life simple. The intention
is to keep up with the latest standards in this space. In support of this, we will provide some form
of fingerprinting for the coded model so that consumers can identify crypto key sets, etc. that are
compatible with this configuration. This fingerprint may simply take the form of some condensed, but
readable packing of relevant bits such as hash algorithm, etc.

*/


type AsymmetricIfc interface {
}

type asymmetric struct {
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewAsymmetric() *asymmetric {
	return &asymmetric{}
}


// -------------------------------------------------------------------------------------------------
// AsymmetricIfc
// -------------------------------------------------------------------------------------------------

// ref: https://medium.com/@Raulgzm/golang-cryptography-rsa-asymmetric-algorithm-e91363a2f7b3

// TODO: Our simple form needs to support a single public/private key pair so that one application
// can create/sign/store, and another application can read (without the ability to create/sign store
// Optionally account for and/or add support for two-party exchanges, each with their own key pairs
// in support of encrypted communications


