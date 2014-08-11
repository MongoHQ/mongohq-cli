# MongoHQ CLI

For usage instructions see the [documentation](http://docs.compose.io/getting-started/cli.html).

This is a test MongoHQ cli project.  It works using the API and Gopher
projects.  The purpose of this project is to:

* Allow customers to access features in a CLI that would give them more functionality.  Such as:
  * Log file tailing in the CLI would enable using pipes and advanced functionality like `logs | grep`.
  * Stats monitoring from the command line, which is much nicer than
* Allow MongoHQ developers to give access to features without the need for UI features.

## Installation

To install, run:

```
curl https://compose-cli.s3.amazonaws.com/install.sh | sh
```

## Using in Dev Mode

```
git clone git@github.com:MongoHQ/mongohq-cli.git
cd mongohq-cli
go run *.go deployments
```

## Files

* `mongohq.go` is a router for commands
* `*_controller.go` are the controllers and views based on the data returned from the api.
* `*_api.go` are the methods for interacting with the API

## Conventions

Each set of actions has its own `controller` and corresponding `api`.  The following conventions work:

* Routers only call the controller action
* Controller actions call the API to get data, the API does not call the controller.
* APIs return errors to the controllers, and the controllers will print the error.

## MVP

* Authenticate a user with 2fa (if enabled) (complete)
* List databases (complete)
* List deployments (complete)
* Tail mongostat (complete)
* Query database log (complete)
