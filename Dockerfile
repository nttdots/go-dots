FROM centos:7

ENV GOPATH /go
ENV PATH /usr/local/go/bin:${PATH}
RUN yum install -y git make gcc gnutls-devel nc which && \
    mkdir go && \
    cd /usr/local/ &&  \
    curl https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz > go1.8.1.linux-amd64.tar.gz && \
    tar xf go1.8.1.linux-amd64.tar.gz && rm go1.8.1.linux-amd64.tar.gz &&\
    curl https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh > /usr/bin/wait-for-it.sh && \
    chmod 755  /usr/bin/wait-for-it.sh

ENV dots_dir /go/src/github.com/nttdots/go-dots

ENV PATH /go/bin:${PATH}

RUN go get "github.com/nttdots/go-dots/..."

ADD https://api.github.com/repos/nttdots/go-dots/git/refs/heads/master .

RUN go get "github.com/nttdots/go-dots/..."

RUN chmod 755 ${dots_dir}/dots_client/entry_point.sh && \
    chmod 755 ${dots_dir}/dots_server/entry_point.sh