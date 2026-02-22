FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /meter-simulator ./cmd/simulator

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Shanghai

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /meter-simulator /app/meter-simulator

RUN chown -R appuser:appuser /app

USER appuser

CMD ["/app/meter-simulator"]
