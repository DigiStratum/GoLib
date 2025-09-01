package oauth2

/*

Standard Oauth2 Access Token implementation using GUID style tokens as an AccessTokenIfc

TODO:
 * Add a TokenTime property and support for adjustable timesource in both the builder and the JSON
   parser
*/

import (
	gojson "encoding/json"
	"strings"

	chrono "github.com/DigiStratum/GoLib/Chrono"
)

type GuidAccessTokenIfc interface {
	// Embedded interface(s)
	AccessTokenIfc

	// Our own interface
}

type guidAccessToken struct {
	// Standard OAuth2 fields
	AccessToken  string `json:"access_token"`            // Required
	TokenType    string `json:"token_type"`              // Required
	ExpiresIn    int64  `json:"expires_in,omitempty"`    // Recommended
	RefreshToken string `json:"refresh_token,omitempty"` // Optional
	Scopes       string `json:"scope,omitempty"`         // Optional
	TokenId      string `json:"id_token,omitempty"`      // Optional

	expiresAt *chrono.TimeStamp
	scopeList []string
}

// -------------------------------------------------------------------------------------------------
// Builder
// -------------------------------------------------------------------------------------------------

type guidAccessTokenBuilder struct {
	guidAccessToken
}

func NewGuidAccessTokenBuilder() *guidAccessTokenBuilder {
	// TODO: leverage a token store to generate a unique token (GUID style vs JWT signed)
	return &guidAccessTokenBuilder{
		guidAccessToken: guidAccessToken{
			TokenType: "Bearer",
		},
	}
}

func (b *guidAccessTokenBuilder) SetAccessToken(token string) *guidAccessTokenBuilder {
	b.AccessToken = token
	return b
}

func (b *guidAccessTokenBuilder) SetRefreshToken(token string) *guidAccessTokenBuilder {
	b.RefreshToken = token
	return b
}

func (b *guidAccessTokenBuilder) SetExpiresIn(seconds int64) *guidAccessTokenBuilder {
	if seconds > 0 {
		b.expiresAt = NewTokenTime().ExpiresAt((seconds))
		b.ExpiresIn = seconds
	}
	return b
}

func (b *guidAccessTokenBuilder) SetScopes(scopes []string) *guidAccessTokenBuilder {
	b.scopeList = scopes
	b.Scopes = ""
	if len(scopes) > 0 {
		b.Scopes = strings.Join(scopes, " ")
	}
	return b
}

func (b *guidAccessTokenBuilder) SetTokenId(tokenId string) *guidAccessTokenBuilder {
	b.TokenId = tokenId
	return b
}

func (b *guidAccessTokenBuilder) Build() *guidAccessToken {
	return &b.guidAccessToken
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

/*
	{
	  "access_token": "YOUR_ACCESS_TOKEN_STRING",
	  "token_type": "Bearer",
	  "expires_in": 3600,
	  "refresh_token": "YOUR_REFRESH_TOKEN_STRING",
	  "scope": "scope1 scope2 scope3",
	  "id_token": "YOUR_ID_TOKEN_STRING"
	}
*/

// Convert RFC 6749 Access Token Response JSON into AccessTokenIfc instance
func NewGuidAccessTokenFromJson(jsonData []byte) *guidAccessToken {
	r := guidAccessToken{}
	err := gojson.Unmarshal(jsonData, &r)
	if err != nil {
		return nil
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// AccessTokenIfc
// -------------------------------------------------------------------------------------------------

func (r *guidAccessToken) IsValid() bool {
	// TODO: Implement the logic to check if the access token is valid
	return false
}

// -------------------------------------------------------------------------------------------------
// GoLib/Data/json/JsonSerializableIfc
// -------------------------------------------------------------------------------------------------

func (r *guidAccessToken) ToJson() (*string, error) {
	jsonBytes, err := gojson.Marshal(r)
	if nil != err {
		return nil, err
	}
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}
