FROM golang:1.9.2-alpine as builder

# Install dep
RUN apk add --update ca-certificates git && \
    go get -u github.com/golang/dep/cmd/dep

# Build project
WORKDIR /go/src/github.com/batazor/go-auth
COPY . .
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest

RUN addgroup -S 997 && adduser -S -g 997 997
USER 997

WORKDIR /app/
COPY --from=builder /go/src/github.com/batazor/go-auth/app .
CMD ["./app"]
