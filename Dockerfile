FROM ubuntu:12.04

RUN apt-get update -qq && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  openssl git mercurial subversion bzr

ADD https://storage.googleapis.com/golang/go1.5.linux-amd64.tar.gz /usr/local/go1.5.tar.gz
WORKDIR /usr/local
RUN tar -C /usr/local -xzf go1.5.tar.gz
RUN ln -s /usr/local/go/bin/go /usr/local/bin/go

RUN adduser --gecos '' --disabled-password harvest

ENV GOPATH /home/harvest

RUN install -o harvest -d /home/harvest/src/github.com/mitch000001/go-harvest

ADD . /home/harvest/src/github.com/mitch000001/go-harvest

WORKDIR /home/harvest/src/github.com/mitch000001/go-harvest

RUN git remote set-url origin https://github.com/mitch000001/go-harvest.git

RUN go get -u ./...
