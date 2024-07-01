# **Go Transactions Service**

- [Go Transactions Service](#Go-Transactions-Service)
   - [Description](#Description)
   - [Setup](#setup)
      - [Prerequisites](#prerequisites)
   - [Pre-setup](#pre-setup)
      - [Local Setup](#local-setup)
      - [Docker setup](#docker-setup)
   - [How to test the program](#How-to-test-the-program)
   - [Makefile commands](#makefile-commands)
      - [golangci-lint](#golangci-lint)
      - [govulncheck](#govulncheck)
      - [mockgen](#mockgen)
      - [sqlc](#sqlc)
      - [migrations](#migrations)
   - [Built With](#built-with)

## Description
This repository contains the Go Transactor Service, a solution for processing transaction files(CSV) and sending summary emails. The project features comprehensive file handling, database interactions, and email notifications.
## Setup

### Prerequisites

Ensure that you have the following prerequisites installed:

- [Golang](https://go.dev/) - version 1.22.4
- [Docker](https://docs.docker.com/engine/install/)
- [Make](https://makefiletutorial.com/)

### Pre-setup

1. **Be sure Local Configuration File is correct**:

    - This repository is designed so that you can run this repository locally, however if you need to make any special adjustments to your DB or SMTP server configurations, you can make your changes in file `config.yaml`.
    - Your configuration should look similar or same to the example. Example:

      ```yaml
      ENVIRONMENT_NAME: local
      DB_TYPE: "postgresql"
      DB_HOST: localhost
      DB_PORT: "5432"
      DB_USERNAME: root
      DB_PASSWORD: ""
      DB_SCHEMA: gostori
      SMTP_HOST: "localhost"
      SMTP_PORT: "2525"
      MAX_BATCH_SIZE: 1000
      ```

### Local Setup

Run the following commands:

```bash
# Update Go module dependencies
go mod tidy

# Verify code formatting, linting, and security/vulnerability risks
make verify

# Run tests
make test

# In case you need to update your golang version on your machine
make update_go

# Run Docker Compose
make docker-up

# Run the migrations files to the DB
make migrate-up

# Run the main program
make run file=path/to/your/file.csv email=your_email@email.com
```

### Docker Setup
If you are done with the previous steps you can omit this part.
```bash
# Update Go module dependencies
go mod tidy

# Verify code formatting, linting, and security/vulnerability risks
make verify

# Run tests
make test

# Run Docker Compose
make docker-up
```

### How to test the program
Make sure you have the path to the csv file you want to test.
If you don't have one in the `cvs-example` folder you can find one (`txns-csv`) to do your tests with more than 1400 transactions.
Make sure that you are in the root of the repository and then run the following command to execute the program:

```bash
# Run the main program
make run file=path/to/your/file.csv email=your_email@email.com
```
If everything went well with the program you will see in console a message like "Transactions processed successfully".
Then you can go to the following address on your computer, and you should be able to see the email you received:
http://localhost:3000/ `SMTP Container`

![See Email Example](https://i.imgur.com/Zan7tqv.png)

### sqlc

`sqlc` is a tool for Go that generates type-safe Go code from SQL queries. (Use this commands for development purposes)

```bash
make sqlc
```

### migrations

`migrate` is a tool for Go that generates type-safe Go code from SQL queries.

```bash
make NAME="$FILE_NAME" migrate-new
make migrate-up
make migrate-down
make migrate-force
```

## Built With

- [Go](https://go.dev/) - version 1.22.4
- [Sqlc](https://sqlc.dev/)
- [zap](go.uber.org/zap)
- [Makefile](https://www.gnu.org/software/make/manual/make.html#Introduction)
- [Testify](https://github.com/stretchr/testify)
- [Mockgen](go.uber.org/mock/mockgen)
- [Dockertest](github.com/ory/dockertest)