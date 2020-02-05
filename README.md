# Calendar API

Creating an API for an interview calendar

<!-- vim-markdown-toc GFM -->

* [Prerequisites](#prerequisites)
* [Run the app](#run-the-app)

<!-- vim-markdown-toc -->

## Prerequisites

* [go](https://golang.org)
* [gocov](https://github.com/axw/gocov#installation) (required for the `make test`
  command)


## Run the app

It should be possible to run the app using the make commands defined in the
Makefile.

To create a new database, run:

    make db-create

To delete the database, run:

    make db-drop

To run the server, run:

    make serve

To run the test suite:

    make test
