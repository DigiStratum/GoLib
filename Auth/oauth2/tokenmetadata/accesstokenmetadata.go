package tokenmetadata

/*

TODO:
 * Implement builder pattern to support immutable result for accessTokenMetadata

*/
import (
	"encoding/json"
)

type AccessTokenMetadataIfc interface {
	GetIssuedAt() int64
	SetIssuedAt(issuedAt int64)
	GetExpiresAt() int64
	SetExpiresAt(expiresAt int64)
	GetScope() string
	SetScope(scope string)
}

type accessTokenMetadata struct {
	IssuedAt  int64  `json:"issued"`  // UTC Timestamp when this token was issued
	ExpiresAt int64  `json:"expires"` // UTC Timestamp when this token expires
	Scope     string `json:"scope"`   // Scope(s) attached to this token (from client_details)
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewAccessTokenMetadata() AccessTokenMetadataIfc {
	atm := accessTokenMetadata{}
	return &atm
}

// -------------------------------------------------------------------------------------------------
// AccessTokenMetadataIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *accessTokenMetadata) GetIssuedAt() int64 {
	return r.IssuedAt
}

func (r *accessTokenMetadata) SetIssuedAt(issuedAt int64) {
	r.IssuedAt = issuedAt
}

func (r *accessTokenMetadata) GetExpiresAt() int64 {
	return r.ExpiresAt
}

func (r *accessTokenMetadata) SetExpiresAt(expiresAt int64) {
	r.ExpiresAt = expiresAt
}

func (r *accessTokenMetadata) GetScope() string {
	return r.Scope
}

func (r *accessTokenMetadata) SetScope(scope string) {
	r.Scope = scope
}

// -------------------------------------------------------------------------------------------------
// JsonSerializableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *accessTokenMetadata) ToJson() (*string, error) {
	jsonBytes, err := json.Marshal(r)
	if nil != err {
		return nil, err
	}
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}

// -------------------------------------------------------------------------------------------------
// JsonDeserializableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *accessTokenMetadata) FromJson(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), r)
}
