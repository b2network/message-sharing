FROM golang:latest as builder
#ARG GOPROXY=https://goproxy.cn
COPY . /src
RUN apt-get update && \
    cd /src/applications && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /src/build/builder cmd/builder/main.go && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /src/build/listener cmd/listener/main.go && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /src/build/proposer cmd/proposer/main.go && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /src/build/validator cmd/validator/main.go
# RUN apt-get update && apt install -y protobuf-compiler git && \
#    cd /tmp && git clone https://github.com/googleapis/googleapis.git && \
#    cp -r /tmp/googleapis/* /usr/local/include/ && \
#    cd /src && \
#    go get github.com/golang/protobuf/protoc-gen-go && \
#    go install github.com/golang/protobuf/protoc-gen-go && \
#    PROTO_INCLUDE=/usr/local/include make proto && \
#    GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/main cmd/app/main.go


# =====
FROM alpine:latest
WORKDIR "/src"
# bin
COPY --from=builder /src/build/builder /usr/bin/builder
COPY --from=builder /src/build/listener /usr/bin/listener
COPY --from=builder /src/build/proposer /usr/bin/proposer
COPY --from=builder /src/build/validator /usr/bin/validator
# config
COPY --from=builder /src/applications/config/builder.yaml /src/config/builder.yaml
COPY --from=builder /src/applications/config/listener.yaml /src/config/listener.yaml
COPY --from=builder /src/applications/config/proposer.yaml /src/config/proposer.yaml
COPY --from=builder /src/applications/config/validator.yaml /src/config/validator.yaml

CMD ["/usr/bin/builder -f /src/config/builder.yaml"]
