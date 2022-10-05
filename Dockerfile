FROM golang:1.17.2-buster

ENV ENVIRONMENT=production

WORKDIR /usr/app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build .

CMD [ "./polarite" ]
