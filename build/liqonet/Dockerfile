FROM golang:1.22.2 as goBuilder-wg

ARG VERSION=2e0774f246fb4fc1bd5cb44584d033038c89174e
# change with "go install git.zx2c4.com/wireguard-go"
# waiting for https://github.com/WireGuard/wireguard-go/pull/87 to be merged
RUN git clone https://git.zx2c4.com/wireguard-go
WORKDIR /go/wireguard-go
RUN git checkout $VERSION
RUN CGO_ENABLED=0 make

FROM golang:1.22.2 as goBuilder
WORKDIR /tmp/builder

COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN  go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$(go env GOARCH) go build -ldflags="-s -w" ./cmd/liqonet


FROM alpine:3.19

RUN apk update && \
    apk add iptables ipset bash wireguard-tools tcpdump conntrack-tools curl && \
    rm -rf /var/cache/apk/*

COPY --from=goBuilder /tmp/builder/liqonet /usr/bin/liqonet
COPY --from=goBuilder-wg /go/wireguard-go/wireguard-go /usr/bin/wireguard-go

ENTRYPOINT [ "/usr/bin/liqonet" ]
