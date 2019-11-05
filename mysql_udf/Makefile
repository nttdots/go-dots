UNAME_S=$(shell uname -s)
ifeq ($(UNAME_S),Darwin)
INC_MYSQL=-I/usr/local/include/mysql
else
INC_MYSQL=-I$(shell dirname $(shell find / -name mysql.h 2>/dev/null))
endif

mysql_udf:
	gcc -c -Wall -fpic mysql-notification.c $(INC_MYSQL)
	gcc -shared -o mysql-notification.so mysql-notification.o
install:

clean:
	rm -f mysql-notification.o
	rm -f mysql-notification.so
