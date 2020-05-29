FROM ubuntu:18.04 AS builder
COPY sources.list /etc/apt/sources.list
RUN apt-get update && apt-get install -y git
WORKDIR /src
RUN git clone https://github.com/juxuny/bridge.git


FROM golang:1.12.5
WORKDIR /go/src/github.com/juxuny
COPY --from=builder /src/bridge ./bridge
RUN go install github.com/juxuny/bridge/cmd/bridge-server
WORKDIR /app
ENTRYPOINT ["bridge-server", "-c=server.json", "-v=true"]