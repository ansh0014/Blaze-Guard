FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/citizenalert-agent ./agent/citizenalert

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /out/citizenalert-agent /app/citizenalert-agent

EXPOSE 8004

CMD ["/app/citizenalert-agent"]
