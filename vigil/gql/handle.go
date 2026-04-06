package gql

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/twiglab/h2o/vigil/gql/graph"
)

func Handle(conf graph.Config) http.Handler {
	srv := handler.New(graph.NewExecutableSchema(conf))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	r := chi.NewRouter()

	r.Handle("/", playground.ApolloSandboxHandler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	return r
}
