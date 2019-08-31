BIN_NAME=adcl

build:
	go build -o ${BIN_NAME}

install:
	go build -o ${GOPATH}/bin/${BIN_NAME}

uninstall:
	rm -f ${GOPATH}/bin/${BIN_NAME}

clean:
	rm -f ${BIN_NAME}