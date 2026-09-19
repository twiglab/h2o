package gql

import (
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/nab/gql/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

type conf struct {
	path string
}

type Ops func(*conf)

func WithPath(p string) Ops {
	return func(c *conf) {
		c.path = p
	}
}

func Handle(hd nab.HandleData, ops ...Ops) http.Handler {
	c := &conf{path: "gql"}

	for _, o := range ops {
		o(c)
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		HandleData: hd,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	mux := chi.NewMux()

	mux.Handle("/", playground.ApolloSandboxHandler("GraphQL playground", path.Join(c.path, "/query")))
	mux.Handle("/query", srv)

	return srv
}
