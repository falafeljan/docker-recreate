SOURCEDIR=.
SOURCES := $(find $(SOURCEDIR) -name '*.go')

BINARY=docker-recreate
LDFLAGS=-ldflags "-X main.BuildTime=`date +%FT%T%z`"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
		go build ${LDFLAGS} -o ${BINARY} ${SOURCES}

.PHONY: install
install:
		go install ${LDFLAGS} ./...

.PHONY: clean
clean:
		if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
