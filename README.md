# Calendar API

Creating an API for an interview calendar. This was originally intended to be
part of a coding challenge for a role I applied for but I had issues with my
laptop and fucked up. So now I'm using this to try out ideas for a basic CRUD
application in Golang.

<!-- vim-markdown-toc GFM -->

* [Prerequisites](#prerequisites)
* [Install](#install)
* [Run the app](#run-the-app)

<!-- vim-markdown-toc -->

## Prerequisites

* [go](https://golang.org)
* [PostgreSQL v10+](https://www.postgresql.org/)
* [gocov](https://github.com/axw/gocov#installation) (required for the `make test`
  command)
* [go-swagger](https://github.com/go-swagger/go-swagger/)
* [swagger-merger](https://github.com/WindomZ/swagger-merger)
  * The library above is an NPM package so this would require
    [Node/NPM](https://nodejs.org/en/) as well

## Install

To install the app run:

    go get -u github.com/kohrVid/calendar-api

## Run the app

It should be possible to run the app using the make commands defined in the
Makefile.

To create a new database, run:

    make db-create

To delete the database, run:

    make db-drop

To run migrations:

    make db-migrate

To reverse a recent migration:

    make db-migrate-down

To run the server, run:

    make serve

To view the swagger documentation, run:

    make swagger

To run the test suite:

    make test

Run the test suite with inotify:

    make test-hot-reload

All of these commands should work in the test environment (`ENV=test`). For commands
such as `make db-seed`, it may be necessary to update the config/env.yaml file
before running in the development environment.
