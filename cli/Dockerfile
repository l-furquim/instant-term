FROM golang:1.24.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY server/*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

CMD ["/docker-gs-ping"]