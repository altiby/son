FROM golang:1.20-alpine3.17 as builder
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -mod vendor -o "cli" cmd/cli/main.go

FROM scratch
COPY --from=builder /app/cli .
COPY  configs configs
COPY  migrations migrations
COPY resources/people.csv.zip people.csv.zip
CMD ["./cli", "init", "users", "-file=./people.csv.zip"]
