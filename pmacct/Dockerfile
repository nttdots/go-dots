FROM ubuntu:16.04

ENV PMACCT_VERSION=1.6.2
ENV PMACCT_SRC_TAR=pmacct-$PMACCT_VERSION.tar.gz
ENV PMACCT_SRC_DIR=/usr/local/src/pmacct-$PMACCT_VERSION

WORKDIR /usr/local/src

RUN apt-get update && \
    apt-get install -y make gcc pkg-config libtool autoconf automake \
    libmysqlclient-dev libpcap-dev wget && \
    rm -rf /var/lib/apt/lists/* && \
    wget http://www.pmacct.net/$PMACCT_SRC_TAR -O /usr/local/src/$PMACCT_SRC_TAR && \
    tar xf /usr/local/src/$PMACCT_SRC_TAR && \
    rm $PMACCT_SRC_TAR && \
    cd $PMACCT_SRC_DIR && \
    ./autogen.sh && \
    ./configure --enable-mysql && \
    make -j 4 && \
    make install && \
    rm -r $PMACCT_SRC_DIR

CMD nfacctd -f /etc/nfacctd.conf