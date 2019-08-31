BIN_NAME=adcl

all: build windows mac

build:
	go build -o ${BIN_NAME}

windows:
	GOOS=windows GOARCH=386 go build -o ${BIN_NAME}.exe

mac:
	GOOS=darwin GOARCH=amd64 go build -o ${BIN_NAME}_darwin

install:
	go build -o ${GOPATH}/bin/${BIN_NAME}

uninstall:
	rm -f ${GOPATH}/bin/${BIN_NAME}

clean:
	rm -f ${BIN_NAME}
	rm -f ${BIN_NAME}_darwin
	rm -f ${BIN_NAME}.exe