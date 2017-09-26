FROM alpine

RUN apk --update add freeradius

ADD radiusd.conf /etc/raddb/radiusd.conf
ADD clients.conf /etc/raddb/clients.conf
ADD authorize /etc/raddb/mods-config/files/authorize

EXPOSE 1812/udp 1813/udp

CMD ["radiusd", "-X"]
