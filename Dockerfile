FROM golang:1.17.6-alpine3.15 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build main.go

FROM scratch
COPY --from=build /app/main /app
ENTRYPOINT ["/app"]