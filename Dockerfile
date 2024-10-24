FROM golang:1.23-alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.49.0

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

CMD ["air", "-c", ".air.toml"]