UNAME_S=$(shell uname -s)
ifeq ($(UNAME_S),Darwin)
INC_MYSQL=-I/usr/local/include/mysql
PLUGIN_MYSQL=/usr/local/lib/plugin/
MYSQL_RESTART=mysql.server restart
else
INC_MYSQL=-I$(shell dirname $(shell find / -name mysql.h 2>/dev/null))
PLUGIN_MYSQL=/usr/lib/mysql/plugin/.
MYSQL_RESTART=sudo /etc/init.d/mysql restart
endif

mysql_udf:
	gcc -c -Wall -fpic mysql-notification.c $(INC_MYSQL)
	gcc -shared -o mysql-notification.so mysql-notification.o
install:
	sudo cp mysql-notification.so $(PLUGIN_MYSQL)
	$(MYSQL_RESTART)
clean:
	rm -f mysql-notification.o
	rm -f mysql-notification.so
