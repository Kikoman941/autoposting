package resolver

import (
	"autoposting/internal/presentation/graphql/gen"
	"context"
)

func (r *mutationResolver) CreateSocialNetworkAccount(
	ctx context.Context,
	input gen.CreateSocialNetworkAccountInput,
) (gen.CreateSocialNetworkAccountOutput, error) {
	out, err := r.usecase.SocialNetwork.CreateSocialNetworkAccount(ctx, input)
	if err != nil {
		return nil, NewResolverError(
			"Не удалось создать аккаунт",
			err,
		)
	}

	return out, nil
}

func (r *mutationResolver) CreateSocialNetworkPage(
	ctx context.Context,
	input gen.CreateSocialNetworkPageInput,
) (gen.CreateSocialNetworkPageOutput, error) {
	out, err := r.usecase.SocialNetwork.CreateSocialNetworkPage(ctx, input)
	if err != nil {
		return nil, NewResolverError(
			"Не удалось создать страницу соц сети",
			err,
		)
	}

	return out, nil
}

func (r *mutationResolver) CreatePost(
	ctx context.Context,
	input gen.CreatePostInput,
) (gen.CreatePostOutput, error) {
	out, err := r.usecase.SocialNetwork.CreatePost(ctx, input)
	if err != nil {
		return nil, NewResolverError(
			"Ну удалось опубликовать пост",
			err,
		)
	}

	return out, nil
}
