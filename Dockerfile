FROM pingcap/go-rust:latest

RUN mkdir -p /deploy/bin

# tiadmin
COPY . /go/src/github.com/pingcap/tiadmin
RUN cd /go/src/github.com/pingcap/tiadmin/tiadmin \
    && go clean ../... \
    && go-wrapper download \
    && go-wrapper install \
    && cp /go/bin/tiadmin /deploy/bin/tiadmin

# tidb
RUN git clone https://github.com/pingcap/tidb.git /go/src/github.com/pingcap/tidb \
    && cd /go/src/github.com/pingcap/tidb \
    && make parser && make server \
    && cp tidb-server/tidb-server /deploy/bin/tidb-server

# pd
RUN git clone https://github.com/pingcap/pd.git /go/src/github.com/pingcap/pd \
    && cd /go/src/github.com/pingcap/pd \
    && make build \
    && cp bin/pd-server /deploy/bin/pd-server

# tikv
RUN git clone https://github.com/pingcap/tikv.git /go/src/github.com/pingcap/tikv \
    && cd /go/src/github.com/pingcap/tikv \
    && cargo build --release \
    && cp target/release/tikv-server /deploy/bin/tikv-server \
    && cp target/release/tikv-dump /deploy/bin/tikv-dump \
    && cargo clean

# cleanup
RUN rm -rf /go/bin/* /go/pkg/* /go/src/*

WORKDIR /deploy
EXPOSE 80 1234 6060 4000 10080 5551
ENTRYPOINT ["/deploy/bin/tiadmin"]
