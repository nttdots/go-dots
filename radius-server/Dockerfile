FROM alpine

RUN apk --update add freeradius

ADD radiusd.conf /etc/raddb/radiusd.conf
ADD clients.conf /etc/raddb/clients.conf
ADD authorize /etc/raddb/mods-config/files/authorize

CMD ["radiusd", "-X"]
