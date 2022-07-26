# step 1: build base go image
FROM golang:1.18-alpine as builder

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o loggerApp ./cmd/api
RUN chmod +x /app/loggerApp

# step 2: build tiny docker image
FROM alpine:latest
RUN mkdir /app
COPY --from=builder /app/loggerApp /app
CMD ["/app/loggerApp"]
