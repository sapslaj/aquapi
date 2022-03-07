package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sapslaj/aquapi/cmd/aquapi-admin-ui/graph"
	"github.com/sapslaj/aquapi/cmd/aquapi-admin-ui/graph/generated"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	ui := http.FileServer(http.Dir("./ui"))
	http.Handle("/", ui)
	http.Handle("/graphiql", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("UI running at http://localhost:%s/", port)
	log.Printf("GraphiQL running at http://localhost:%s/graphiql", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
