package oauth2

import (
	chrono "github.com/DigiStratum/GoLib/Chrono"
)

type AccessTokenIfc interface {
	IsValid() bool
	GetAccessToken() string
	GetRefreshToken() string

	GetExpiresIn() int64
	GetScope() string
	GetExpiresAt() *chrono.TimeStamp
}

type accessToken struct {
	accessToken  string
	refreshToken string
}

func NewAccessToken() *accessToken {
	// TODO: leverage a token store to generate a unique token (GUID style vs JWT signed)
	return &accessToken{
		accessToken:  "TODO",
		refreshToken: "TODO",
	}
}

func (r *accessToken) IsValid() bool {
	// TODO: Implement the logic to check if the access token is valid
	return false
}
