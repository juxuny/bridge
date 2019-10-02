FROM ubuntu
RUN apt-get update && apt-get install -y git golang-go 
RUN mkdir /gopath
ENV GOPATH=/gopath
RUN go get github.com/juxuny/bridge
RUN go install github.com/juxuny/bridge/cmd/bridge-server
WORKDIR /app
COPY ./server.json /app
COPY ./token.conf /app
ENTRYPOINT ["/gopath/bin/bridge-server", "-c", "server.json"]
