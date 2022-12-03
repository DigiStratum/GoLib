package identityprovider

type IdentityProviderCollectionIfc interface {
	Put(idp IdentityProviderIfc)
	Get(idpId string) IdentityProviderIfc
}

type identityProviderCollection struct {
	idps		map[string]IdentityProviderIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------
func NewIdentityProviderCollection() IdentityProviderCollectionIfc {
	ipc := identityProviderCollection{
		idps:	make(map[string]IdentityProviderIfc),
	}
	return &ipc
}

// -------------------------------------------------------------------------------------------------
// IdentityProviderCollectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (ipc *identityProviderCollection) Put(idp IdentityProviderIfc) {
	(*ipc).idps[idp.GetId()] = idp
}

func (ipc *identityProviderCollection) Get(idpId string) IdentityProviderIfc {
	if idp, ok := (*ipc).idps[idpId]; ok { return idp }
	return nil
}
