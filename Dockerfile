FROM golang:1.20-bookworm AS builder

ENV ENVIRONMENT=production

WORKDIR /usr/app

COPY . .

RUN go build .

FROM debian:bookworm AS runtime

WORKDIR /app

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates openssl curl

COPY --from=builder /usr/app/* .

COPY --from=builder /usr/app/views ./views

RUN mkdir -p /app/data/

ENV HTTP_PORT=3000

ENV DATABASE_DIRECTORY=/app/data/

ENV ENVIRONMENT=production

EXPOSE ${HTTP_PORT}

CMD [ "./polarite" ]
