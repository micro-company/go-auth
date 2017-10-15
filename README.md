# go-auth-server

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
| MONGO_URL        | localhost/auth       |

### technology stack

#### Back-End

* Go
* MongoDB

#### Library

+ [chi](github.com/pressly/chi) - for routing
+ [glide](github.com/Masterminds/glide) - for vendoring