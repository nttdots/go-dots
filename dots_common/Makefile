
bindata.go: $(wildcard schemas/*.json) $(GOPATH)/bin/go-bindata
	$(GOPATH)/bin/go-bindata -o $@ -pkg dots_common schemas/*.json

$(GOPATH)/bin/go-bindata:
	go get github.com/jteeuwen/go-bindata/...

clean:
	rm -f bindata.go

test:
	go test