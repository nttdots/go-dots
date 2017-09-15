FROM alpine
RUN apk update \
 && apk add build-base cmake wget tar gzip file linux-headers bison flex glib libgcrypt glib-dev libidn-dev gnutls-dev libgcrypt-dev gnutls libidn gnutls-utils ca-certificates lksctp-tools lksctp-tools-dev\
 && wget -O freediameter.tar.gz http://www.freediameter.net/hg/freeDiameter/archive/1.2.1.tar.gz \
 && mkdir build-freediameter \
 && cd build-freediameter \
 && tar xzvf ../freediameter.tar.gz --strip-components=1 \
 && cmake -DCMAKE_INSTALL_PREFIX=/usr -DDEFAULT_CONF_PATH=/etc/freeDiameter -DBUILD_TESTING:BOOL=OFF -Wno-dev .\
 && make \
 && make install \
 && cd - \
 && rm -R freediameter.tar.gz build-freediameter \
 && apk del build-base cmake wget tar gzip linux-headers glib-dev bison flex libgcrypt-dev libidn-dev gnutls-dev git

COPY freeDiameter.conf /
ADD ["start.sh", "/"]

EXPOSE 3868 5658
CMD ["/start.sh"]
