FROM golang:latest
RUN apt-get update && apt-get -y install \
    automake \
    bison \
    flex \
    gcc \
    libjansson-dev \
    libmagic-dev \
    libssl-dev \
    libtool \
    make \
    wget

# Config & Build Yara
WORKDIR /develop
RUN wget https://github.com/VirusTotal/yara/archive/v3.10.0.tar.gz
RUN tar -zxf v3.10.0.tar.gz
WORKDIR /develop/yara-3.10.0
RUN ./bootstrap.sh
RUN ./configure --enable-cuckoo --enable-magic --enable-dotnet --with-crypto
RUN make
RUN make install
RUN make check

# build YaraPerfTest
WORKDIR /develop
COPY go.mod /develop/
COPY *.go /develop/
RUN go build -o YaraPerfTest


EXPOSE 1234
ENTRYPOINT [ "YaraPerfTest" ]