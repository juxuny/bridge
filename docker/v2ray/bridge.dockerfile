FROM ubuntu:18.04 AS builder
RUN apt-get update && apt-get install -y git
WORKDIR /src
RUN git clone https://github.com/juxuny/bridge.git


FROM golang:1.12.5
WORKDIR /go/src/github.com/juxuny
COPY --from=builder /src/bridge ./bridge
RUN go install github.com/juxuny/bridge/cmd/bridge-client
WORKDIR /app
ENTRYPOINT ["bridge-client", "-c=client.json", "-v=true", "-t=20"]
