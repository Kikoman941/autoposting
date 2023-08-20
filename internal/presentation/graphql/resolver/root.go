package resolver

import (
	"autoposting/internal/app/registry"
	"autoposting/internal/presentation/graphql/gen"
)

type Resolver struct {
	usecase *registry.Usecases
}

func NewResolver(u *registry.Usecases) *Resolver {
	return &Resolver{
		u,
	}
}

func (r *Resolver) Mutation() gen.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() gen.QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
