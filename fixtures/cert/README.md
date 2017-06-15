# Client Certificates

I used [these instructions][howto] to generate the client certificates in this directory.

Very roughly:

``` sh
# Generate a root (enter a password of your choice)
openssl genrsa -des3 -out ca.key 4096
openssl req -new -x509 -days 7000 -key ca.key -out ca.crt

# Generate a client key and certificate (enter a
# password of your choice)
openssl genrsa -des3 -out client.key 4096
openssl req -new -key client.key -out client.csr
openssl x509 -req -days 7000 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out client.crt

# Combine client key/cert into an unprotected PKCS
# (when prompted, enter an empty password)
openssl pkcs12 -export -clcerts -in client.crt -inkey client.key -out fixtures/cert/client.p12
```

[howto]: https://gist.github.com/mtigas/952344
