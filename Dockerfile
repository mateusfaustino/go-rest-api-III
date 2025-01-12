FROM golang:latest

WORKDIR /app

# CMD [ "tail", "-f", "/dev/null" ]

COPY . .

RUN go mod tidy

EXPOSE 8080

CMD [ "go", "run", "cmd/server/main.go" ]