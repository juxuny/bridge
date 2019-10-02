FROM ubuntu
RUN apt-get update && apt-get install -y git golang-go 
RUN mkdir /gopath
ENV GOPATH=/gopath
RUN go get github.com/juxuny/bridge
RUN go install github.com/juxuny/bridge/cmd/bridge-client
WORKDIR /app
COPY ./client.json /app
ENTRYPOINT ["/gopath/bin/bridge-client", "-c", "client.json", "-v=false"]
