NAME = website-verification

VERSION ?= v0.0.1

DIST_FOLDER := dist

RELEASE_FOLDER := resources


.PHONY: build container clean
build: 
	go mod tidy
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ./${DIST_FOLDER}/${NAME} ./cmd/${NAME}


container: build
	docker build -t ${NAME}:${VERSION} -f ${RELEASE_FOLDER}/Dockerfile .;

