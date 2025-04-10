FROM golang:1.24-bookworm AS builder

ENV ENVIRONMENT=production

WORKDIR /usr/app

COPY . .

RUN go build .

FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates openssl curl

COPY . .

COPY --from=builder /usr/app/polarite /usr/local/bin/polarite

RUN mkdir -p /app/data/

ENV HTTP_PORT=3000

ENV DATABASE_DIRECTORY=/app/data/

ENV ENVIRONMENT=production

EXPOSE ${HTTP_PORT}

CMD [ "/usr/local/bin/polarite" ]
