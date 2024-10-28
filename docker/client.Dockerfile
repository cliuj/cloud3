FROM golang:1.23.2

WORKDIR /opt/client

COPY . ./

RUN cd cmd/client && go build -o client
RUN mkdir -p /tmp/cloud3/client

ENTRYPOINT ["./cmd/client/client"]
