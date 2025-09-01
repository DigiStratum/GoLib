# Oauth2 support for https://rfc-editor.org/rfc/rfc6749.html
----
Note: maybe abstract https://pkg.go.dev/golang.org/x/oauth2 with better conveniences
----
or leverage: https://github.com/go-oauth2/oauth2/blob/master/example/server/server.go
https://github.com/go-oauth2/oauth2/tree/master/example
https://medium.com/@bytecraze.com/create-an-oauth2-server-in-15-minutes-using-go-a660f6246e61
----

A successful OAuth 2.0 JSON access token API response typically follows a standard format defined in the OAuth 2.0 specification (RFC 6749) and related specifications like RFC 6750 for Bearer tokens. The response is a JSON object and includes the following key parameters:

{
  "access_token": "YOUR_ACCESS_TOKEN_STRING",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "YOUR_REFRESH_TOKEN_STRING", 
  "scope": "openid email profile",
  "id_token": "YOUR_ID_TOKEN_STRING" 
}

Explanation of Parameters:
access_token (REQUIRED):
This is the actual token that the client application uses to authenticate requests to the protected resource (e.g., an API). Its format can vary, but it's commonly a randomly generated string or a JSON Web Token (JWT).
token_type (REQUIRED):
Indicates the type of token issued. For most OAuth 2.0 implementations, this value is "Bearer", as described in RFC 6750.
expires_in (RECOMMENDED):
The lifetime in seconds of the access token. The client should use this to determine when the token will expire and potentially request a new one using a refresh token.
refresh_token (OPTIONAL):
A token used to obtain new access tokens without requiring the user to re-authenticate. This is typically included in responses to authorization code or password grant types.
scope (OPTIONAL):
The scope of access granted by the access token. This reflects the permissions the client application has been granted.
id_token (OPTIONAL):
A JSON Web Token (JWT) that contains information about the authenticated end-user. This is primarily used in OpenID Connect, an identity layer built on top of OAuth 2.0.
Important Considerations:
HTTP Headers:
The response should also include the Content-Type: application/json header and Cache-Control: no-store to prevent caching of the sensitive token information.
Error Responses:
If the access token request is invalid or fails, the server should return an appropriate error response, typically with an HTTP status code like 400 (Bad Request) and a JSON body detailing the error.
