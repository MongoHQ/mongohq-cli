# MongoHQ CLI

This is a test MongoHQ cli project.  It works using the API and Gopher
projects.  The purpose of this project is to:

* Allow customers to access features in a CLI that would give them more functionality.  Such as:
  * Log file tailing in the CLI would enable using pipes and advanced functionality like `logs | grep`.
  * Stats monitoring from the command line, which is much nicer than
* Allow MongoHQ developers to give access to features without the need for UI features.

## Installation 

To install, run:

```
curl https://mongohq-cli.s3.amazonaws.com/install.sh | sh
```

## Using in Dev Mode

```
git clone git@github.com:MongoHQ/mongohq-cli.git
cd mongohq-cli
source .env
git submodule foreach git pull
go run mongohq.go
```

## Files

* `mongohq.go` is a router for commands
* `src/github.com/MongoHQ/controllers/*` are the controllers and views based on the data returned from the api.
* `src/github.com/MongoHQ/api/*` are the methods fo interacting with the API

## Conventions

Each set of actions has its own `controller` and corresponding `api`.  The following conversations work:

* Routers only call the controller action
* Controller actions call the API to get data, the API does not call the controller.
* APIs return errors to the controllers, and the controllers will print the error.
* Favor explicit typing with `var` instead of implicit typing with `:=` for simplicity.

## MVP

* Authenticate a user with 2fa (if enabled)
* List databases
* List deployments
* Tail mongostat
* Tail log
