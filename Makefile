MAIN_PACKAGE := github.com/NaufalA/wmb-graphql-server/cmd/server

run:	## Run Server application
	@go run "${MAIN_PACKAGE}"
gqlgenerate:
	@go get github.com/99designs/gqlgen && go run github.com/99designs/gqlgen generate