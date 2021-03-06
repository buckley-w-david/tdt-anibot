FROM golang:1.11

WORKDIR /opt/anibot
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/bot ./cmd/bot

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /opt/anibot/bin/bot .
CMD ["./bot"]
