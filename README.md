# go-auth-server

[![Maintainability](https://api.codeclimate.com/v1/badges/32dbb31cde6ea8f52cf0/maintainability)](https://codeclimate.com/github/batazor/go-auth/maintainability)

Auth micro-service

### RUN

```
docker-compose build
docker-compose up
```

### ENV

| Name ENV         | Default value             |
|------------------|---------------------------|
| PORT             | 4070                      |
| MONGO_URL        | localhost/auth            |

### Generation cert

```
openssl genrsa \
    -passout pass:12345678 \
    -out cert/private_key.pem \
    2048
    
openssl rsa \
    -passin pass:12345678 \
    -in cert/private_key.pem \
    -pubout > cert/public_key.pub
```

### technology stack

#### Back-End

* Go
* MongoDB

#### Library

+ [chi](github.com/pressly/chi) - for routing
+ [glide](github.com/Masterminds/glide) - for vendoring
