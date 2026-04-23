package graph

import "github.com/twiglab/h2o/vigil/orm/ent"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

const (
	MaxStr = "~"
	MinStr = "!"
)

type Resolver struct {
	Client *ent.Client
}
