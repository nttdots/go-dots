.PHONY:

dots_client_controller:	$(wildcard ./*.go) $(wildcard ./**/*.go) $(wildcard ../dots_common/**/*.go) \
                $(wildcard ../coap/*.go) ../dots_common/bindata.go ../dots_common/*.go
	go build

../dots_common/bindata.go:
	make -C ../dots_common bindata.go

install:	dots_client_controller
	go install

clean:
	make -C ../dots_common clean
	rm -f dots_client_controller

test:
	go test