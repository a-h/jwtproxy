# Reverse Proxy with JWT Authentication

A reverse proxy which rejects incoming requests which:

* Don't have an `authorization` HTTP header.
* Have an `authorization` header which doesn't contain a JWT.
* Have an `authorization` header which contains an expired or invalid JWT.
* Have an `authorization` header which contains a JWT which has an unrecognised issuer.
* Have an `authorization` header which could not be validated using the public key corresponding to the issuer.

## Usage

* Decide on an issuer (`iss`) value to use for each API client.
   * This is usually a domain, e.g. `example.com`.
* Get the API consumer to generate a private RSA key, and send you the public key.
* Setup the proxy to allow requests from the API consumer using environment variables, or command line flags.
* Get the API consumer to send HTTP requests which set an `authorization` header containing a JWT signed with the private key.

## Generating RSA keys

Generate a private key. (example_private.pem).

```bash
openssl genrsa -out example_private.pem 2048
```

Extract the public key. (example_public.pem)

```bash
openssl rsa -in example_private.pem -outform PEM -pubout -out example_public.pem
```

You start the proxy passing it a map of issuers to public keys as environment variables or a file.

## Minimal JWT

The JWT passed by the client must meet the following criteria:

### Header

The JWT must be signed using the RS256 algorithm.

```json
{
  "alg": "RS256",
  "typ": "JWT"
}
```

### Payload

The payload must contain an issuer agreed between the two parties, and an expiration timestamp.

```json
{
  "iss": "custom_issuer.example.com",
  "exp": "1504282460"
}
```

## Configuration

Configuration can be provided by command line flags (specified by `-name`) or by environment variables.

### JWTPROXY_REMOTE_URL / -remoteURL

The URL to proxy requests to.

### JWTPROXY_REMOTE_HOST_HEADER / -remoteHostHeader

The HTTP host header to send to the remote endpoint (useful if the remote endpoint is not using DNS).

### JWTPROXY_LISTEN_PORT / -port

The TCP port to open up the proxy on.

### JWTPROXY_HEALTHCHECK_URI / -health

The location that the proxy should use to respond to health check HTTP requests (defaults to `/health`).

### JWTPROXY_PREFIX / -prefix

The prefix to strip from incoming requests applied to the remote URL, e.g to make incoming HTTP request `/api/user?id=1` map to outgoing HTTP request `/user?id=1`

### JWTPROXY_ISSUER_ / JWTPROXY_PUBLIC_KEY_

It's possible to set issuer to public key maps by using environment variables alone.

To set the issuer "example.com" to a public key, create two environment variables with a matching suffix for the `JWTPROXY_ISSUER_` and `JWTPROXY_PUBLIC_KEY_` keys, e.g.:

```
JWTPROXY_ISSUER_0=example.com
JWTPROXY_PUBLIC_KEY_0=-----BEGIN PUBLIC KEY-----MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxFj26fqmulXntc7kCp9tMs6MEQUsk2r16Jd6k+aZSaLBo0dVgP77q1os10gZT4N0gYH6NsbVqP4+wWAUIDiemhpxq986z5mtB/lGvmHmaQcK/bOnEvcLWinHJZIla1m2RF7diN5/WBRNh8CyYMiW+BV/6dngknBtP7bDpnCkYrySaOQtKRvrech1UFRKgQjD8bprrcUmOFWYrmKe2NCxcQs9RhYuACt3Du2Z4VwVWN2xvL5LlZdWK7jLENe3MkOZU5WcwA7n+K/tulqA9uNRv8cRIL/y8BUwUsUoqBiyVZXQUa7BgE82GoTXtv3uqkN/yZxnlEcaJW5BD1nFzuvuyQIDAQAB-----END PUBLIC KEY-----
```

### JWTPR# Reverse Proxy with JWT Authentication

A reverse proxy which rejects incoming requests which:

* Don't have an `authorization` HTTP header.
* Have an `authorization` header which doesn't contain a JWT.
* Have an `authorization` header which contains an expired or invalid JWT.
* Have an `authorization` header which contains a JWT which has an unrecognised issuer.
* Have an `authorization` header which could not be validated using the public key corresponding to the issuer.

## Usage

* Decide on an issuer (`iss`) value to use for each API client.
  * This is usually a domain, e.g. `example.com`.
* Get the API consumer to generate a private RSA key, and send you the public key.
* Setup the proxy to allow requests from the API consumer using environment variables, or command line flags.
* Get the API consumer to send HTTP requests which set an `authorization` header containing a JWT signed with the private key.

## Generating RSA keys

Generate a private key. (example_private.pem).

```bash
openssl genrsa -out example_private.pem 2048
```

Extract the public key. (example_public.pem)

```bash
openssl rsa -in example_private.pem -outform PEM -pubout -out example_public.pem
```

You start the proxy passing it a map of issuers to public keys as environment variables or a file.

## Minimal JWT

The JWT passed by the client must meet the following criteria:

### Header

The JWT must be signed using the RS256 algorithm.

```json
{
  "alg": "RS256",
  "typ": "JWT"
}
```

### Payload

The payload must contain an issuer agreed between the two parties, and an expiration timestamp.

```json
{
  "iss": "custom_issuer.example.com",
  "exp": "1504282460"
}
```

## Configuration

Configuration can be provided by command line flags (specified by `-name`) or by environment variables.

### JWTPROXY_REMOTE_URL / -remoteURL

The URL to proxy requests to.

### JWTPROXY_REMOTE_HOST_HEADER / -remoteHostHeader

The HTTP host header to send to the remote endpoint (useful if the remote endpoint is not using DNS).

### JWTPROXY_LISTEN_PORT / -port

The TCP port to open up the proxy on.

### JWTPROXY_HEALTHCHECK_URI / -health

The location that the proxy should use to respond to health check HTTP requests (defaults to `/health`).

### JWTPROXY_PREFIX / -prefix

The prefix to strip from incoming requests applied to the remote URL, e.g to make incoming HTTP request `/api/user?id=1` map to outgoing HTTP request `/user?id=1`

### JWTPROXY_ISSUER_ / JWTPROXY_PUBLIC_KEY_

It's possible to set issuer to public key maps by using environment variables alone.

To set the issuer "example.com" to a public key, create two environment variables with a matching suffix for the `JWTPROXY_ISSUER_` and `JWTPROXY_PUBLIC_KEY_` keys, e.g.:

```
JWTPROXY_ISSUER_0=example.com
JWTPROXY_PUBLIC_KEY_0=-----BEGIN PUBLIC KEY-----MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxFj26fqmulXntc7kCp9tMs6MEQUsk2r16Jd6k+aZSaLBo0dVgP77q1os10gZT4N0gYH6NsbVqP4+wWAUIDiemhpxq986z5mtB/lGvmHmaQcK/bOnEvcLWinHJZIla1m2RF7diN5/WBRNh8CyYMiW+BV/6dngknBtP7bDpnCkYrySaOQtKRvrech1UFRKgQjD8bprrcUmOFWYrmKe2NCxcQs9RhYuACt3Du2Z4VwVWN2xvL5LlZdWK7jLENe3MkOZU5WcwA7n+K/tulqA9uNRv8cRIL/y8BUwUsUoqBiyVZXQUa7BgE82GoTXtv3uqkN/yZxnlEcaJW5BD1nFzuvuyQIDAQAB-----END PUBLIC KEY-----
```

### JWTPROXY_CONFIG / -keys

The location of a JSON file containing a map of issuers to public keys, e.g.:

```bash
JWTPROXY_CONFIG=keys.json
```

JSON does not support newlines, so the public key value will need to have them replaced with \n in JSON, e.g. by `cat dev_pub.pem | tr '\n' '_' | sed 's/_/\\n/g' > dev_pub2.pem`

- keys.json
```json
{
    "example.com": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxFj26fqmulXntc7kCp9t\nMs6MEQUsk2r16Jd6k+aZSaLBo0dVgP77q1os10gZT4N0gYH6NsbVqP4+wWAUIDie\nmhpxq986z5mtB/lGvmHmaQcK/bOnEvcLWinHJZIla1m2RF7diN5/WBRNh8CyYMiW\n+BV/6dngknBtP7bDpnCkYrySaOQtKRvrech1UFRKgQjD8bprrcUmOFWYrmKe2NCx\ncQs9RhYuACt3Du2Z4VwVWN2xvL5LlZdWK7jLENe3MkOZU5WcwA7n+K/tulqA9uNR\nv8cRIL/y8BUwUsUoqBiyVZXQUa7BgE82GoTXtv3uqkN/yZxnlEcaJW5BD1nFzuvu\nyQIDAQAB\n-----END PUBLIC KEY-----"
}
```

### JWTPROXY_STRICT_SSL / -strictSsl
Set strict SSL/TLS off or on by setting the env or passing command with a boolean.
False sets strict ssl off. 
True turns strict ssl on, this is the default value.

# Running it

## Command line

```bash
jwtproxy -remoteURL http://example.com:8080 -keys keys.json
```

## Docker

Or you can run the Docker container, using environment variables to pass in required data. In this case, exposing the linked container 'hopeful_pike'.

```bash
docker run --link hopeful_pike -p 9090:9090/tcp --rm -e "JWTPROXY_LISTEN_PORT=9090" -e "JWTPROXY_REMOTE_URL=http://hopeful_pike:8080" -v /users/me/keys.json:/keys.json -e "JWTPROXY_CONFIG=keys.json" adrianhesketh/jwtproxy
```

# Usage

Use `curl` to access your local proxy.

```bash
curl -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOiIxNDg2MzkyMjAwIiwiZXhwIjoiMTU4NjM5MjIwMCIsImlzcyI6ImV4YW1wbGUuY29tIn0.McHwUVbe96y-vaTOExPjANm8e8p0v6I7puPf74SV7Jn-QYprrhLlklnBP4MEF77v0LIBUFKgzpOMfldCONId3ktOFOf0117x9iWG3J-Zf6Ni3HinhA9U1pPU7_OhTtkXacmgats8tLWAqmOz46NeyAmHS_dkvodUUPpcHY-AqQtzM4ql6RZpMDstz5dFJWZh9P0_prPknoI-argt2jn-KGajCOIghcGxNarylq5oX62rT9AavavyWGnJW0zLnP9qtIuChzigU542Nbg7y6_E7FaVA2cPICPuiPehn6vVTKuVil0o2SJgFD2J2HQfxa0iDrc8HzbubMGJcw7Vlpkk0w" http://localhost:9090
```

# Testing

Generate a JWT with an appropriate payload at [jwt.io], or using a library:

* iat
  * issued at time: The time when the JWT was generated as a Unix timestamp.
* exp
  * expiry time: The time when the JWT expires, may be rejected by the server if the difference between exp and iat is too long.
* iss
  * the issuer, used to look up the correct public key to use to validate the JWT signature.

```json
{
  "iat": "1486392200",
  "exp": "1586392200",
  "iss": "example.com"
}
```

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOiIxNDg2MzkyMjAwIiwiZXhwIjoiMTU4NjM5MjIwMCIsIm5hbWUiOiJBZHJpYW4gSGVza2V0aCJ9.d45Or2h-lApJ4FK2pKj0ZIRdDTULsNl1z-V3LfQFvno
```

To sign and validate using the command line (to compare against the Go implementation):

```bash
# SHA256 hash the data.
openssl dgst -sha256 -binary data.json > hash.bin
# base64 encode the hash so that it should match the value of the X-Sha256hash HTTP header.
openssl base64 -e -in hash.bin -out hash.b64

# Sign the hash.
openssl rsautl -in hash.bin -inkey private_test.pem -sign -out signature.bin
# base64 encode the signature so that it should match the value of the X-Signature HTTP header.
openssl base64 -e -in signature.bin -out signature.b64

# Verify the signature.
openssl rsautl -in signature.bin -verify -inkey public_test.pem -pubin > verified.bin
openssl base64 -e -in verified.bin -out verified.b64

# Compare the original hash to the hash created by the verification routine.
# The two files should be equal.
cat hash.b64
cat verified.b64
```
OXY_CONFIG / -keys

The location of a JSON file containing a map of issuers to public keys, e.g.:

```bash
JWTPROXY_CONFIG=keys.json
```

JSON does not support newlines, so the public key value will need to have them replaced with \n in JSON, e.g. by `cat dev_pub.pem | tr '\n' '_' | sed 's/_/\\n/g' > dev_pub2.pem`

- keys.json
```json
{
    "example.com": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxFj26fqmulXntc7kCp9t\nMs6MEQUsk2r16Jd6k+aZSaLBo0dVgP77q1os10gZT4N0gYH6NsbVqP4+wWAUIDie\nmhpxq986z5mtB/lGvmHmaQcK/bOnEvcLWinHJZIla1m2RF7diN5/WBRNh8CyYMiW\n+BV/6dngknBtP7bDpnCkYrySaOQtKRvrech1UFRKgQjD8bprrcUmOFWYrmKe2NCx\ncQs9RhYuACt3Du2Z4VwVWN2xvL5LlZdWK7jLENe3MkOZU5WcwA7n+K/tulqA9uNR\nv8cRIL/y8BUwUsUoqBiyVZXQUa7BgE82GoTXtv3uqkN/yZxnlEcaJW5BD1nFzuvu\nyQIDAQAB\n-----END PUBLIC KEY-----"
}
```

# Running it

## Command line

```bash
jwtproxy -remoteURL http://example.com:8080 -keys keys.json
```

## Docker

Or you can run the Docker container, using environment variables to pass in required data. In this case, exposing the linked container 'hopeful_pike'.

```bash
docker run --link hopeful_pike -p 9090:9090/tcp --rm -e "JWTPROXY_LISTEN_PORT=9090" -e "JWTPROXY_REMOTE_URL=http://hopeful_pike:8080" -v /users/me/keys.json:/keys.json -e "JWTPROXY_CONFIG=keys.json" adrianhesketh/jwtproxy
```

# Usage

Use `curl` to access your local proxy.

```bash
curl -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOiIxNDg2MzkyMjAwIiwiZXhwIjoiMTU4NjM5MjIwMCIsImlzcyI6ImV4YW1wbGUuY29tIn0.McHwUVbe96y-vaTOExPjANm8e8p0v6I7puPf74SV7Jn-QYprrhLlklnBP4MEF77v0LIBUFKgzpOMfldCONId3ktOFOf0117x9iWG3J-Zf6Ni3HinhA9U1pPU7_OhTtkXacmgats8tLWAqmOz46NeyAmHS_dkvodUUPpcHY-AqQtzM4ql6RZpMDstz5dFJWZh9P0_prPknoI-argt2jn-KGajCOIghcGxNarylq5oX62rT9AavavyWGnJW0zLnP9qtIuChzigU542Nbg7y6_E7FaVA2cPICPuiPehn6vVTKuVil0o2SJgFD2J2HQfxa0iDrc8HzbubMGJcw7Vlpkk0w" http://localhost:9090
```

# Testing

Generate a JWT with an appropriate payload at [jwt.io], or using a library:

* iat
  * issued at time: The time when the JWT was generated as a Unix timestamp.
* exp
  * expiry time: The time when the JWT expires, may be rejected by the server if the difference between exp and iat is too long.
* iss
  * the issuer, used to look up the correct public key to use to validate the JWT signature.

```json
{
  "iat": "1486392200",
  "exp": "1586392200",
  "iss": "example.com"
}
```

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOiIxNDg2MzkyMjAwIiwiZXhwIjoiMTU4NjM5MjIwMCIsIm5hbWUiOiJBZHJpYW4gSGVza2V0aCJ9.d45Or2h-lApJ4FK2pKj0ZIRdDTULsNl1z-V3LfQFvno
```

To sign and validate using the command line (to compare against the Go implementation):

```bash
# SHA256 hash the data.
openssl dgst -sha256 -binary data.json > hash.bin
# base64 encode the hash so that it should match the value of the X-Sha256hash HTTP header.
openssl base64 -e -in hash.bin -out hash.b64

# Sign the hash.
openssl rsautl -in hash.bin -inkey private_test.pem -sign -out signature.bin
# base64 encode the signature so that it should match the value of the X-Signature HTTP header.
openssl base64 -e -in signature.bin -out signature.b64

# Verify the signature.
openssl rsautl -in signature.bin -verify -inkey public_test.pem -pubin > verified.bin
openssl base64 -e -in verified.bin -out verified.b64

# Compare the original hash to the hash created by the verification routine.
# The two files should be equal.
cat hash.b64
cat verified.b64
```
