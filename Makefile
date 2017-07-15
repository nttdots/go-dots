.PHONY: all

all:
	$(CROSS_FLAGS) make -C dots_server/
	$(CROSS_FLAGS) make -C dots_client/
	$(CROSS_FLAGS) make -C dots_client_controller/

install:
	$(CROSS_FLAGS) make -C dots_server/ install
	$(CROSS_FLAGS) make -C dots_client/ install
	$(CROSS_FLAGS) make -C dots_client_controller/ install

test:
	$(CROSS_FLAGS) make -C dots_server/ test
	$(CROSS_FLAGS) make -C dots_client/ test
	$(CROSS_FLAGS) make -C dots_client_controller/ test
clean:
	make -C dots_client clean
	make -C dots_server clean
	make -C dots_client_controller clean