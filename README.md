# web-analyzer

Web-application for analyzing Webpages built in Golang.

## Objective

The objective is to build a web application that does an analysis of a
web-page/URL. The application shows a form with a text field in which users can type in the URL of the webpage to be analyzed. Additionally to the form, it contains a button to send a request to the server.

More details about the requirements can be found in the [requirement document](docs/requirement.md).

## Structure

The project is structured as follows:

```tree
.
├── cmd
│   └── server
│       └── main.go
├── internal
│   └── handlers
│       └── handlers.go
├── docs
│   └── requirement.md
├── README.md
└── go.mod

```

For more details, see also the [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

## Build

run `make build` to build the application.

## Run

run `make run` to start the application. By default, the server will start on port 8080. You can access it by navigating to `http://localhost:8080` in your web browser.

## Test

run `make test` to execute the unit tests for the application.

## Deploy

## Docker

run `make docker-build` to build the Docker image for the application.
run `make docker-run` to run the Docker container for the application.
run `make docker-stop` to stop the Docker container for the application.

## Improvements
