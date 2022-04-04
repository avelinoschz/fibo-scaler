FROM golang:1.17.6-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o bin/main 

FROM scratch
COPY --from=builder /app/bin/main /
ENTRYPOINT ["/main"]