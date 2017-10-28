# go-auth-server

[![Maintainability](https://api.codeclimate.com/v1/badges/32dbb31cde6ea8f52cf0/maintainability)](https://codeclimate.com/github/batazor/go-auth/maintainability)
[![Docker Build Statu](https://img.shields.io/docker/build/jrottenberg/ffmpeg.svg)](https://hub.docker.com/r/batazor/go-auth)

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
