BINARY_NAME=kidex
default: ${BINARY_NAME}

${BINARY_NAME}: *.go
	@go build -o ${BINARY_NAME}

clean:
	@rm ${BINARY_NAME}
