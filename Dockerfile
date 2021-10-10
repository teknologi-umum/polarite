FROM golang:1.17.1-buster

ENV ENVIRONMENT=production

WORKDIR /usr/app

COPY . .

RUN go mod download

RUN go build main.go app.go

CMD [ "./main" ]