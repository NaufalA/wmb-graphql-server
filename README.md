# WMB GraphQL Go Server

GraphQL Go API Implementation using 99designs/gqlgen and gin with mongoDB.

## How to run development server
this example is using docker to run the mongoDB server. just cha
- run `docker compose up -d` to start mongoDB server and mongo express UI (default on port 8081)
- run `go mod download` to resolve Go dependencies
- rename file `.env.example` to `.env` to define environment variables with default values
- run `make run` to start Go server. if you're on vscode run Launch Package config on the Debug menu to start Go server in debug mode
- you can initialize some dummy data and default user credentials by running `go run cmd/seeder/main.go` or the Run Seeder config. User credentials will be printed to terminal when it first created.