FROM golang:1.23.2

WORKDIR /opt/server

COPY . ./

RUN cd cmd/server && go build -o server

ENTRYPOINT ["./cmd/server/server"]
