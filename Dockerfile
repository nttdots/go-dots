FROM centos:7

ENV GOPATH /go
ENV PATH /usr/local/go/bin:/go/bin:${PATH}

ARG BRANCH
ENV BRANCH=${BRANCH:-master}
ENV DOTS_DIR /go/src/github.com/nttdots/go-dots

ADD https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz /root
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /usr/bin

RUN yum install -y git make gcc gnutls-devel nc which && \
    yum clean all && \
    tar xf /root/go1.8.1.linux-amd64.tar.gz -C /usr/local && \
    rm /root/go1.8.1.linux-amd64.tar.gz && \
    chmod 755 /usr/bin/wait-for-it.sh

#    mkdir go && \
#    cd /usr/local/ &&  \
#    curl https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz > go1.8.1.linux-amd64.tar.gz && \
#    tar xf go1.8.1.linux-amd64.tar.gz && rm go1.8.1.linux-amd64.tar.gz &&\
#    curl https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh > /usr/bin/wait-for-it.sh && \
#    chmod 755  /usr/bin/wait-for-it.sh


RUN go get "layeh.com/radius" && \
    go get "github.com/nttdots/go-dots/..."

ADD https://api.github.com/repos/nttdots/go-dots/git/refs/heads/${BRANCH} ${DOTS_DIR}

RUN cd ${DOTS_DIR} && \
    git fetch origin ${BRANCH} && \
    if [ `git rev-parse --abbrev-ref HEAD` == ${BRANCH} ]; then \
          git checkout ${BRANCH}; \
    fi; \
    git merge origin/${BRANCH} && \
    make install

RUN chmod 755 ${DOTS_DIR}/dots_client/entry_point.sh && \
    chmod 755 ${DOTS_DIR}/dots_server/entry_point.sh

EXPOSE 4646 4647
