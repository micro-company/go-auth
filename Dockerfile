FROM alpine:latest

RUN addgroup -S 997 && adduser -S -g 997 997
USER 997

WORKDIR /app/
COPY app app
CMD ["./app"]