VERSION=0.1.1
BIN_NAME=adcl

all: build windows mac

build:
	go build -o ${BIN_NAME} ./cmd/

windows:
	GOOS=windows GOARCH=386 go build -o ${BIN_NAME}.exe ./cmd/

mac:
	GOOS=darwin GOARCH=amd64 go build -o ${BIN_NAME}_darwin ./cmd/

releases: all
	tar -czvf ${BIN_NAME}v${VERSION}_linux.tar.gz ${BIN_NAME}
	tar -czvf ${BIN_NAME}v${VERSION}_win.tar.gz ${BIN_NAME}.exe
	tar -czvf ${BIN_NAME}v${VERSION}_darwin.tar.gz ${BIN_NAME}_darwin

install:
	go build -o ${GOPATH}/bin/${BIN_NAME} ./cmd/

uninstall:
	rm -f ${GOPATH}/bin/${BIN_NAME}

clean_bin:
	rm -f ${BIN_NAME}
	rm -f ${BIN_NAME}_darwin
	rm -f ${BIN_NAME}.exe

clean_releases:
	rm -f ${BIN_NAME}v${VERSION}_linux.tar.gz
	rm -f ${BIN_NAME}v${VERSION}_win.tar.gz
	rm -f ${BIN_NAME}v${VERSION}_darwin.tar.gz

clean: clean_bin clean_releases
