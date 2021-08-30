package transcoder

type EncodingScheme int

// ref: https://pkg.go.dev/encoding#pkg-functions
const (
	ES_UNKNOWN EncodingScheme = iota
	ES_AUTO				// Automagically detect Encoding
	ES_NONE				// No Encoding
	ES_BASE64			// Base 64 Encoding
	ES_UUENCODE			// UU-Encoding (EMAIL)
	ES_HTTPESCAPE			// HTTP Escaped Encoding (HTTP/URL/form-post)
)
