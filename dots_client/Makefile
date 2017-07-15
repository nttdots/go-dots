.PHONY:


dots_client:	$(wildcard ./*.go) $(wildcard ./**/*.go) $(wildcard ../dots_common/**/*.go) \
                $(wildcard ../coap/*.go) ../dots_common/bindata.go
	go build

../dots_common/bindata.go:
	make -C ../dots_common bindata.go

install:	dots_client
	go install

clean:
	make -C ../dots_common clean
	rm -f dots_client
	rm -f dots_client.sock

test:
	go test
