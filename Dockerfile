FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o blogService cmd/app/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/blogService .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./blogService"]