## Reverse Proxy with JWT Authentication

A reverse proxy which validates incoming requests based on the contents of the Authorization header.

## Usage

You provide your API consumer with a unique issuer value (usually a domain, e.g. example.com) and your
API consumer generates a private RSA key, and sends you the public key.

```
openssl genrsa -out example_private.pem 2048
openssl rsa -in example_private.pem -outform PEM -pubout -out example_public.pem
```

You start the proxy passing it a map of issuers to public keys. You'll need to replace newlines
with \n in JSON, e.g. by `cat dev_pub.pem | tr '\n' '_' | sed 's/_/\\n/g' > dev_pub2.pem`

- keys.json
```
{
    "example.com": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxFj26fqmulXntc7kCp9t\nMs6MEQUsk2r16Jd6k+aZSaLBo0dVgP77q1os10gZT4N0gYH6NsbVqP4+wWAUIDie\nmhpxq986z5mtB/lGvmHmaQcK/bOnEvcLWinHJZIla1m2RF7diN5/WBRNh8CyYMiW\n+BV/6dngknBtP7bDpnCkYrySaOQtKRvrech1UFRKgQjD8bprrcUmOFWYrmKe2NCx\ncQs9RhYuACt3Du2Z4VwVWN2xvL5LlZdWK7jLENe3MkOZU5WcwA7n+K/tulqA9uNR\nv8cRIL/y8BUwUsUoqBiyVZXQUa7BgE82GoTXtv3uqkN/yZxnlEcaJW5BD1nFzuvu\nyQIDAQAB\n-----END PUBLIC KEY-----"
}
```

You can run it yourself...

```
jwtproxy -remoteURL http://example.com:8080 -keys keys.json
```

Or you can run the Docker container, using environment variables to pass in required data. In this case, exposing the linked container 'hopeful_pike'.

```
docker run --link hopeful_pike -p 9090:9090/tcp --rm -e "JWTPROXY_LISTEN_PORT=9090" -e "JWTPROXY_REMOTE_URL=http://hopeful_pike:8080" -e "JWTPROXY_CONFIG=keys.json" adrianhesketh/jwtproxy
```

Now attempt to use `curl` to access your local proxy.

```
curl -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOiIxNDg2MzkyMjAwIiwiZXhwIjoiMTU4NjM5MjIwMCIsImlzcyI6ImV4YW1wbGUuY29tIn0.McHwUVbe96y-vaTOExPjANm8e8p0v6I7puPf74SV7Jn-QYprrhLlklnBP4MEF77v0LIBUFKgzpOMfldCONId3ktOFOf0117x9iWG3J-Zf6Ni3HinhA9U1pPU7_OhTtkXacmgats8tLWAqmOz46NeyAmHS_dkvodUUPpcHY-AqQtzM4ql6RZpMDstz5dFJWZh9P0_prPknoI-argt2jn-KGajCOIghcGxNarylq5oX62rT9AavavyWGnJW0zLnP9qtIuChzigU542Nbg7y6_E7FaVA2cPICPuiPehn6vVTKuVil0o2SJgFD2J2HQfxa0iDrc8HzbubMGJcw7Vlpkk0w" http://localhost:9090
```

## Testing

Generate a JWT with an appropriate payload at [jwt.io], or using a library:

* iat
  * issued at time: The time when the JWT was generated as a Unix timestamp.
* exp
  * expiry time: The time when the JWT expires, may be rejected by the server if the difference between exp and iat is too long.
* iss
  * the issuer, used to look up the correct public key to use to validate the JWT signature.

```
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

```
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
