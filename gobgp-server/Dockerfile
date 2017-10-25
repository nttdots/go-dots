FROM tokatsu/quagga:0.1

ENV GOPATH /go
ENV PATH /go/bin:$PATH
WORKDIR /root
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git gcc wget && \
    wget https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz && \
    tar xf go1.8.linux-amd64.tar.gz  -C /usr/local/ && rm go1.8.linux-amd64.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go&&  \
    go get github.com/osrg/gobgp/...


COPY start.sh /root
COPY gobgpd.conf /root

CMD ["bash", "/root/start.sh"]

