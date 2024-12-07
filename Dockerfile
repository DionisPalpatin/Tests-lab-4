FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 587/tcp

CMD ["go", "run", "./cmd/main.go"]