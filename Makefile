all: prelog cmdtmpl

prelog:
	go build ./cmd/prelog

cmdtmpl:
	go build ./cmd/cmdtmpl

clean:
	rm -f prelog
	rm -f cmdtmpl

test:
	go test .
