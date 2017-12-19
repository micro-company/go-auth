APP?=go-mail-crd
GRPC_PORT?=50051

clean:
	rm -f ${APP}

build: clean
	go build -o ${APP}

run: build
	GRPC_PORT=${GRPC_PORT} ./${APP}
