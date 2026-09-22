package gql

import (
	"net/http"
	"path"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/twiglab/h2o/chrgg/gql/graph"
	"github.com/twiglab/h2o/chrgg/orm"
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

func Handle(db *orm.DBx, ops ...Ops) http.Handler {

	c := &conf{path: "gql"}

	for _, o := range ops {
		o(c)
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		DBx: db,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	p := path.Join(c.path, "/query")

	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)
	mux.Handle("/", playground.ApolloSandboxHandler("GraphQL playground", p))
	mux.Handle("/query", srv)

	return mux
}
