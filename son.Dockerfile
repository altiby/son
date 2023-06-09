FROM golang:1.20-alpine3.17 as builder
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -mod vendor -o "son" cmd/son/main.go

FROM scratch
COPY --from=builder /app/son .
COPY  configs configs
COPY  migrations migrations
CMD ["./son"]
