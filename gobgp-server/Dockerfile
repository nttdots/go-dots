FROM tokatsu/quagga:0.1

ENV GOPATH /go
ENV PATH /go/bin:$PATH
WORKDIR /root

RUN apt-get update && \
    apt-get install -y git gcc && apt-get -y upgrade ca-certificates && \
    apt-get clean

ADD https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz /root

RUN tar xf /root/go1.8.linux-amd64.tar.gz -C /usr/local && \
    rm -f /root/go1.8.linux-amd64.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go && \
    go get github.com/osrg/gobgp/...

COPY start.sh /root
COPY gobgpd.conf /root

RUN chmod 755 /etc/quagga/init_conf.sh

CMD ["bash", "/root/start.sh"]
