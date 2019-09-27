FROM ubuntu:18.04

RUN apt-get update \
 && apt-get install -y \
    autoconf \
    automake \
    bison \
    build-essential \
    git \
    libtool \
    swig \
    pkg-config \
    python3-dev \
    wget \
 && rm -rf /var/lib/apt/lists/*

RUN cd /tmp \
 && git clone https://github.com/cmusphinx/sphinxbase \
 && cd /tmp/sphinxbase \
 && PYTHON=python3.6 ./autogen.sh LDFLAGS="-L/usr/bin/python3" \
 && make \
 && make install \
 && rm -rf /tmp/sphinxbase \
 && cd /tmp \
 && git clone https://github.com/cmusphinx/pocketsphinx \
 && cd /tmp/pocketsphinx \
 && PYTHON=python3.6 ./autogen.sh LDFLAGS="-L/usr/bin/python3" \
 && make \
 && make install \
 && rm -rf /tmp/pocketsphinx

ENV LD_LIBRARY_PATH="/usr/local/lib"

RUN cd /tmp \
 && wget https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz \
 && tar -C /usr/local -xzf go1.13.1.linux-amd64.tar.gz \
 && rm go1.13.1.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go test
