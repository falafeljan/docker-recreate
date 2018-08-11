SOURCEDIR=.n
SOURCES := $(find $(SOURCEDIR) -name '*.go')

BINARY=build/docker-recreate
LDFLAGS=-ldflags "-X main.BuildTime=`date +%FT%T%z`"

.DEFAULT_GOAL: $(BINARY)

all: clean prebuild test build

.PHONY: prebuild
prebuild: $(SOURCES)
	dep ensure

.PHONY: build
build: $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} ./cli

.PHONY: install
install:
	go install ${LDFLAGS} ./cli

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: test
test:
	go test ./...
