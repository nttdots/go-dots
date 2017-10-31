.PHONY: clean

dots_server: $(wildcard ./*.go) $(wildcard ./**/*.go) $(wildcard ../dots_common/**/*.go) \
             $(wildcard ../coap/*.go) ../dots_common/bindata.go $(wildcard ../dots_common/*.go)
	go build -v

../dots_common/bindata.go:
	make -C ../dots_common bindata.go

install:	dots_server
	go install

clean:
	rm -f dots_server
	make -C ../dots_common clean

test:
	go test
	$(MAKE) -C db_models test
	$(MAKE) -C models test
	$(MAKE) -C controllers test
	$(MAKE) -C radius test

