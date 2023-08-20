package resolver

import (
	"autoposting/internal/presentation/graphql/gen"
	"context"
	"fmt"
)

func (r *queryResolver) GetAccountAuthURL(
	ctx context.Context,
	input gen.GetAccountAuthURLInput,
) (gen.GetAccountAuthURLOutput, error) {
	out, err := r.usecase.SocialNetwork.GetAuthURL(ctx, input)
	if err != nil {
		return nil, NewResolverError(
			fmt.Sprintf("Cannot get auth url for %s", input.SocialNetwork),
			err,
		)
	}
	return out, nil
}

func (r *queryResolver) GetPagesFromSocialNetwork(
	ctx context.Context,
	input gen.GetPagesFromSocialNetworkInput,
) (gen.GetPagesFromSocialNetworkOutput, error) {
	out, err := r.usecase.SocialNetwork.GetPagesFromSocialNetwork(ctx, input)
	if err != nil {
		return nil, NewResolverError(
			fmt.Sprintf("Cannot get account pages for %s", input.SocialNetwork),
			err,
		)
	}
	return out, nil
}
