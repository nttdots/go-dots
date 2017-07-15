
bindata.go: $(wildcard schemas/*.json)
	go get github.com/jteeuwen/go-bindata/...
	$(GOPATH)/bin/go-bindata -o $@ -pkg dots_common schemas/*.json

clean:
	rm -f bindata.go

test:
	go test