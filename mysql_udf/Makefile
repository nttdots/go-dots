mysql_udf:
	gcc -c -Wall -fpic mysql-notification.c -I$(shell dirname $(shell find / -name mysql.h 2>/dev/null))
	gcc -shared -o mysql-notification.so mysql-notification.o
install:
	sudo cp mysql-notification.so /usr/lib/mysql/plugin/.
clean:
	rm -f mysql-notification.o
	rm -f mysql-notification.so
