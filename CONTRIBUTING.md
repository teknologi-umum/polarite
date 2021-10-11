# Contributing Guide

Hello! We'd love to see your contribution on this repository soon, even if it's just a typo fix!

Contributing means anything from reporting bugs, ideas, suggestion, code fix, even new feature.

Bear in mind to keep your contributions under the [Code of Conduct](./github/CODE_OF_CONDUCT.md) for the community.

## Bug report, ideas, and suggestion

The [issues](https://github.com/teknologi-umum/polarite/issues) page is a great way to communicate to us. Other than that, we have a [Telegram group](https://t.me/teknologi_umum) that you can discuss your ideas into. If you're not an Indonesian speaker, it's 100% fine to talk in English there.

Please make sure that the issue you're creating is in as much detail as possible. Poor communication might lead to a big mistake, we're trying to avoid that.

## Pull request

**A big heads up before you're writing a breaking change code or a new feature: Please open up an [issue](https://github.com/teknologi-umum/polarite/issues) regarding what you're working on, or just talk in the [Telegram group](https://t.me/teknologi_umum).**

### Prerequisites

You will need a few things to get things working:

1. Latest version of Go (as of now, we're using v1.17.2). You can get it [here](https://golang.org/dl). For Linux users, you can follow the commands below:

```sh
$ wget https://golang.org/dl/go1.17.2.linux-amd64.tar.gz

$ sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.17.2.linux-amd64.tar.gz

$ export PATH=$PATH:/usr/local/go/bin
# Or put it on your ~/.bashrc or equivalent file.
```

2. [Docker](https://docs.docker.com/engine/install/) and [Docker Compose](https://docs.docker.com/compose/install/) if you want to test the app using docker. (Optional)

### Getting Started

1. [Fork](https://help.github.com/articles/fork-a-repo/) this repository to your own Github account and [clone](https://help.github.com/articles/cloning-a-repository/) it to your local machine.
2. Run `go mod download` to install the dependencies needed.
3. Get the database up and running. You can use `docker-compose up -d` for this. To stop the container, use `docker-compose stop`. To remove the container, use `docker-compose down`. Bear in mind that the data stored in the Redis and MongoDB of the Docker container is not persistent. Once it's stopped, the data will be erased.
4. You can use [Postman](https://www.postman.com/), [Insomnia](https://insomnia.rest/) or [Hoppscotch](https://hoppscotch.io/) to create an API request.
5. Rename `.env.example` to `.env` and fill the config key=value needed. The one's necessary is `ENVIRONMENT`, `PORT`, `DATABASE_URL` and `REDIS_URL`, you may leave everything else blank. If you're using the Docker Compose file to spin up the database, your `.env` should be:

```ini
ENVIRONMENT=development
PORT=3000
DATABASE_URL=mysql://root:password@localhost:3306/polarite
REDIS_URL=redis://@localhost:6379/
LOGTAIL_TOKEN=
SENTRY_DSN=
```

6. Run `go run main.go app.go` to start the development server.
7. Have fun!

You are encouraged to use [Conventional Commit](https://www.conventionalcommits.org/en/v1.0.0-beta.2/) for your commit message. But it's not really compulsory.

### Testing your change

We follow Go convention of testing. The test files would be alongside the original file on the same directory. To run the test, do:

```sh
$ go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
```

```
.
├── business
│  ├── controllers          - Business logic
│  └── models               - Business model
├── CONTRIBUTING.md         - You are here
├── docker-compose.yml      - Docker Compose file
├── Dockerfile
├── go.mod
├── go.sum
├── handlers                - Route handlers and middleware
├── LICENSE
├── main.go                 - Main entrypoint (along with app.go)
├── app.go
├── platform                - Third-party needs to supplement the app
│  ├── highlight
│  ├── logtail
│  ├── migration
│  └── placeholder.sql      - Placeholder data for development
├── README.md
├── repository              - Collections of constants and static variables
└── resources               - Utilities
```

### Before creating a PR

Please test (command above) and run Go FMT to pass the CI.

```sh
$ go fmt ./...
$ go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
```

Please note that we don't enforce some kind of convention and patterns like you would see on Java and other languages.

Go code should be boring. And follow [Go Proverbs](https://go-proverbs.github.io/).

And you're set!
