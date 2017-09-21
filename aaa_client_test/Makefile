
aaa_client_test:	$(wildcard ./*.go) bindata.go
	go build

bindata.go: eap_dict.xml $(GOPATH)/bin/go-bindata
	$(GOPATH)/bin/go-bindata -o $@ -pkg main eap_dict.xml

$(GOPATH)/bin/go-bindata:
	go get github.com/jteeuwen/go-bindata/...

clean:
	rm -f aaa_client_test bindata.go

test:
	go test


