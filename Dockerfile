FROM pingcap/go-rust:latest

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /deploy/bin /deploy/logs /deploy/data

# tiadmin
COPY . /go/src/github.com/pingcap/tiadmin
RUN cd /go/src/github.com/pingcap/tiadmin/tiadmin \
    && go clean ../... \
    && go-wrapper download \
    && go-wrapper install \
    && cp /go/bin/tiadmin /deploy/bin

# tidb
RUN git clone https://github.com/pingcap/tidb.git /go/src/github.com/pingcap/tidb \
	&& git clone https://github.com/golang/snappy.git /go/src/github.com/golang/snappy \
	&& git -C /go/src/github.com/golang/snappy checkout 17e435849f9b5cc7818817bb6e1d5c002f486fef \
    && cd /go/src/github.com/pingcap/tidb \
    && make godep && make parser && make server \
    && cp tidb-server/tidb-server /deploy/bin/tidb-server \
    && make clean

# pd
RUN cd /go/src/github.com/pingcap/pd \
    && make build \
    && cp bin/pd-server /deploy/bin/pd-server

# tikv
RUN git clone https://github.com/pingcap/tikv.git /go/src/github.com/pingcap/tikv \
    && cd /go/src/github.com/pingcap/tikv \
    && cargo build --release \
    && cp target/release/tikv-server /deploy/bin/tikv-server \
    && cargo clean

# cleanup
RUN rm -rf /go/bin/* /go/pkg /go/src/*

WORKDIR /deploy
EXPOSE 80 1234 6060 4000 10080 5551
ENTRYPOINT ["/deploy/bin/tiadmin"]
