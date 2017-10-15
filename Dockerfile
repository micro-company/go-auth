FROM golang:1.9.1-alpine as builder

ENV GLIDE_VERSION v0.13.0

# Install glide
RUN apk add --update ca-certificates wget git && \
    update-ca-certificates && \
    wget https://github.com/Masterminds/glide/releases/download/${GLIDE_VERSION}/glide-${GLIDE_VERSION}-linux-amd64.tar.gz && \
    tar -zxf glide-${GLIDE_VERSION}-linux-amd64.tar.gz && \
    mv linux-amd64/glide /usr/local/bin/

# Build project
WORKDIR /go/src/github.com/batazor/go-bookmarks
COPY . .
RUN glide install
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest

RUN addgroup -S 997 && adduser -S -g 997 997
USER 997

WORKDIR /app/
COPY --from=builder /go/src/github.com/batazor/go-bookmarks/app .
CMD ["./app"]