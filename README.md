# Ublogging

## Uala Challenge

## Requirements

- [Golang 1.23](https://go.dev/doc/devel/release)
- [Docker](https://www.docker.com/)
- [GNU/Make](https://www.gnu.org/software/make/)

## How to run?

```bash
git clone git@github.com:wmb1207/ublogging.git
cd ublogging
go mod tidy

# Set up your .env file (check the wiki for an example)

# Setup the development database
make run-mongo

go run cmd/ublogging/main.go
```

## How to build?

```bash
git clone git@github.com:wmb1207/ublogging.git
cd ublogging
go mod tidy
go build cmd/ublogging/main.go -o ublogging
```