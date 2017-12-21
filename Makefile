test: dep
	dep ensure
	go test .

dep:
ifeq (, $(shell which dep))
	go get -u github.com/golang/dep/cmd/dep
endif

all: prelog cmdtmpl

prelog:
	go build ./cmd/prelog

cmdtmpl:
	go build ./cmd/cmdtmpl

clean:
	rm -f prelog
	rm -f cmdtmpl
