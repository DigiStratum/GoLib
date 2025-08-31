package oauth2

import (
	chrono "github.com/DigiStratum/GoLib/Chrono"
	"github.com/DigiStratum/GoLib/Data/metadata"
)

type AccessTokenType int

const (
	AccessTokenType_Unknown AccessTokenType = iota
	AccessTokenType_Guid
	AccessTokenType_Jwt
)

type AccessTokenIfc interface {
	// Type-specific properties
	GetTokenType() AccessTokenType
	GetMetadata() metadata.MetadataIfc // Scopes, claims, user info, etc

	// Standard properties
	IsValid() bool

	GetAccessToken() string
	GetRefreshToken() string

	GetExpiresIn() int64
	GetExpiresAt() *chrono.TimeStamp
}
