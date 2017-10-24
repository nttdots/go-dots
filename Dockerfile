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

ENV DOTS_DIR /go/src/github.com/nttdots/go-dots

ENV PATH /go/bin:${PATH}

RUN go get "layeh.com/radius"
RUN go get "github.com/nttdots/go-dots/..."

ARG BRANCH
ENV BRANCH=${BRANCH:-master}

ADD https://api.github.com/repos/nttdots/go-dots/git/refs/heads/${BRANCH} .

RUN cd ${DOTS_DIR} && \
    git fetch origin ${BRANCH} && \
    if [ `git rev-parse --abbrev-ref HEAD` == ${BRANCH} ]; then \
          git checkout ${BRANCH}; \
    fi; \
    git merge origin/${BRANCH} && \
    make install

RUN chmod 755 ${DOTS_DIR}/dots_client/entry_point.sh && \
    chmod 755 ${DOTS_DIR}/dots_server/entry_point.sh