package gql

import (
	"github.com/go-chi/chi/v5"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/twiglab/h2o/vigil/gql/graph"
	"github.com/twiglab/h2o/vigil/orm/ent"
)

func Handle(conf graph.Config) *chi.Mux {
	srv := handler.New(graph.NewExecutableSchema(conf))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	r := chi.NewMux()

	r.Handle("/", playground.ApolloSandboxHandler("GraphQL playground", "/gql/query"))
	r.Handle("/query", srv)

	return r
}

func NewConf(cli *ent.Client) graph.Config {
	return graph.Config{
		Resolvers: &graph.Resolver{
			Client: cli,
		},
	}
}
