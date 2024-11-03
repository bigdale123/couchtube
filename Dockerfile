FROM golang:1.22.3-alpine

RUN apk update && apk add --no-cache sqlite

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

COPY static /app/static

EXPOSE 8081

CMD ["./main"]
