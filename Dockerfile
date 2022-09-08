ARG GITHUB_PATH=github.com/Dsmit05/party-day-bot

FROM golang:1.18-alpine AS builder

RUN apk add --no-cache ca-certificates git make

WORKDIR /home/${GITHUB_PATH}

COPY . .

RUN make build

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /home/${GITHUB_PATH}/party-day-bot .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .

EXPOSE 8081
EXPOSE 8082

CMD ["./party-day-bot"]