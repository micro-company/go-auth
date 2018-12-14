# go-auth-server

[![Maintainability](https://api.codeclimate.com/v1/badges/32dbb31cde6ea8f52cf0/maintainability)](https://codeclimate.com/github/batazor/go-auth/maintainability)
[![](https://images.microbadger.com/badges/image/batazor/go-auth.svg)](https://microbadger.com/images/batazor/go-auth "Get your own image badge on microbadger.com")
[![Go Report Card](https://goreportcard.com/badge/github.com/micro-company/go-auth)](https://goreportcard.com/report/github.com/micro-company/go-auth)

Auth micro-service


### Feature

+ JWT auth
+ Support Google recaptcha
+ User manager
+ [UI example](https://micro-company.github.io/react-app/)


![Schema auth-service](docs/schema.png)

### Getting start

```
go get -u github.com/golang/protobuf/proto

protoc -I grpc/mail/ grpc/mail/mail.proto --go_out=plugins=grpc:grpc/mail

docker-compose build
docker-compose up
```

### ENV

| Name ENV                    | Default value                              |
|-----------------------------|--------------------------------------------|
| PORT                        | 4070                                       |
| MONGO_URL                   | mongodb://localhost/auth                   |
| REDIS_URL                   | redis://localhost:6379                     |
| RECAPTCHA_PRIVATE_KEY       | secretKey                                  |
| ENABLE_CAPTCHA              | false                                      |
| **OAuth**                   | --                                         |
| OAUTH_GOOGLE_CLIENT_ID      | --                                         |
| OAUTH_GOOGLE_CLIENT_SECRET  | --                                         |
| OAUTH_REDIRECT_URL          | http://localhost:3000/auth/callback/:type  |

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

### Kubernetes

```
# Run minikube
minikube start \
  --network-plugin=cni \
  --kubernetes-version=v1.8.0
  
# Install Helm
# See https://github.com/kubernetes/helm/blob/master/docs/install.md
helm init
helm repo update

# Run application
helm \
  --kube-context minikube \
  install \
  --name go-auth \
  --namespace=demo \
  ops/Helm/go-auth
  
# Delete
helm del --purge go-auth
```

### Initial state for data base

`initialState/user.json` - contains initial information 

### OAuth

+ [Google setting](https://developers.google.com/identity/protocols/OAuth2)