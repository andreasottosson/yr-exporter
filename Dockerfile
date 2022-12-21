FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./

COPY main.go ./

RUN go build -o /yr-exporter

EXPOSE 9118

CMD ["/yr-exporter"]