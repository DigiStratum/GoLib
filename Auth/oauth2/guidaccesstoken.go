package oauth2

/*

Standard Oauth2 Access Token implementation using GUID style tokens as an AccessTokenIfc

*/

import (
	gojson "encoding/json"

	chrono "github.com/DigiStratum/GoLib/Chrono"
)

type GuidAccessTokenIfc interface {
	AccessTokenIfc
}

type guidAccessToken struct {
	// Standard OAuth2 fields
	AccessToken  string `json:"access_token"`  // The access token issued by the authorization server
	RefreshToken string `json:"refresh_token"` // The refresh token, which can be used to obtain new access tokens using the same authorization grant

	tokenType string
	expiresAt *chrono.TimeStamp
	scopes    []string
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewGuidAccessToken() *guidAccessToken {
	// TODO: leverage a token store to generate a unique token (GUID style vs JWT signed)
	return &guidAccessToken{
		tokenType:    "guid",
		AccessToken:  "TODO",
		RefreshToken: "TODO",
		scopes:       []string{"read"},
	}
}

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
// AccessTokenIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *guidAccessToken) IsValid() bool {
	// TODO: Implement the logic to check if the access token is valid
	return false
}
