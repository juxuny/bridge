FROM golang:1.13.4 AS builder
WORKDIR /pro
COPY . /pro
ENV CGO_ENABLED 0
RUN GOPROXY=https://goproxy.cn go mod download
RUN cd cmd/bridge-server && go build -o /go/bin/bridge-server && cd - 
RUN cd cmd/bridge-client && go build -o /go/bin/bridge-client && cd -

FROM ineva/alpine:3.9
WORKDIR /app
COPY --from=builder /go/bin/bridge-server /app/bridge-server
COPY --from=builder /go/bin/bridge-client /app/bridge-client
ENTRYPOINT /app/bridge-client
