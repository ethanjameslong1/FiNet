FROM golang:1.24

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/finet-server ./FileSearchServer.go

EXPOSE 8080

CMD ["/app/finet-server"]

