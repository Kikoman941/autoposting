package usecase

import (
	"autoposting/internal/domain/service"
	"autoposting/internal/presentation/graphql/gen"
	ewrap "autoposting/pkg/err-wrapper"
	"context"
)

type SocialNetworkUsecase struct {
	socialNetworkService *service.SocialNetworkService
}

func NewSocialNetworkUsecase(
	socialNetworkService *service.SocialNetworkService,
) *SocialNetworkUsecase {
	return &SocialNetworkUsecase{
		socialNetworkService,
	}
}

func (u *SocialNetworkUsecase) CreateSocialNetworkAccount(
	ctx context.Context,
	input gen.CreateSocialNetworkAccountInput,
) (gen.CreateSocialNetworkAccountOutput, error) {
	if err := u.socialNetworkService.CreateSocialNetworkAccount(
		ctx,
		input.SocialNetwork,
		input.Credentials,
	); err != nil {
		switch {
		case service.IsSocialNetworkAccountAlreadyExistsError(err):
			return gen.SocialNetworkAccountAlreadyExistsError{
				Message: err.Error(),
			}, nil
		case service.IsValidationError(err):
			return gen.ValidationError{
				Message: err.Error(),
			}, nil
		case service.IsInternalError(err):
			return gen.InternalError{
				Message: err.Error(),
			}, nil
		default:
			return nil, ewrap.Errorf(
				"failed to create social network %v account with credentials %v",
				input.SocialNetwork,
				input.Credentials,
			)
		}
	}

	return gen.CreateSocialNetworkAccountResult{
		Ok: true,
	}, nil
}

func (u *SocialNetworkUsecase) CreateSocialNetworkPage(
	ctx context.Context,
	input gen.CreateSocialNetworkPageInput,
) (gen.CreateSocialNetworkPageOutput, error) {
	if err := u.socialNetworkService.CreateSocialNetworkPage(ctx, input); err != nil {
		switch {
		case service.IsPageAlreadyExistsError(err):
			return gen.PageAlreadyExistsError{
				Message: err.Error(),
			}, nil
		default:
			return nil, err
		}
	}

	return gen.CreateSocialNetworkPageResult{
		Ok: true,
	}, nil
}

func (u *SocialNetworkUsecase) GetAccessToken(
	ctx context.Context,
	params map[string][]string,
) error {
	if _, ok := params["socialNetwork"]; !ok {
		return ewrap.Errorf("url param socialNetwork not found")
	}
	socialNetwork := params["socialNetwork"][0]
	socialNetworkAccount, err := u.socialNetworkService.GetSocialNetworkAccount(ctx, socialNetwork)
	if err != nil {
		switch {
		case service.IsValidationError(err) || service.IsNotFoundError(err):
			return err
		default:
			return ewrap.Errorf("failed to find social socialNetwork %s account: %w", socialNetwork, err)
		}
	}

	token, err := u.socialNetworkService.GetAccessToken(ctx, socialNetworkAccount, params)
	if err != nil {
		return err
	}

	if err = u.socialNetworkService.SaveAccessToken(ctx, token, socialNetworkAccount); err != nil {
		return err
	}

	return nil
}

func (u *SocialNetworkUsecase) GetAuthURL(
	ctx context.Context,
	input gen.GetAccountAuthURLInput,
) (gen.GetAccountAuthURLOutput, error) {
	socialNetworkAccount, err := u.socialNetworkService.GetSocialNetworkAccount(ctx, input.SocialNetwork)
	if err != nil {
		switch {
		case service.IsValidationError(err):
			return gen.ValidationError{
				Message: err.Error(),
			}, nil
		case service.IsNotFoundError(err):
			return gen.InternalError{
				Message: err.Error(),
			}, nil
		default:
			return nil, ewrap.Errorf(
				"failed to find social socialNetwork %s account: %w", input.SocialNetwork, err,
			)
		}
	}

	authUrl, err := u.socialNetworkService.GetAuthURL(
		socialNetworkAccount.SocialNetwork,
		socialNetworkAccount.Credentials,
	)
	if err != nil {
		if service.IsInternalError(err) {
			return gen.InternalError{
				Message: err.Error(),
			}, nil
		}
		return nil, err
	}

	return gen.GetAccountAuthURLResult{
		URL: authUrl,
	}, nil
}

func (u *SocialNetworkUsecase) GetPagesFromSocialNetwork(
	ctx context.Context,
	input gen.GetPagesFromSocialNetworkInput,
) (gen.GetPagesFromSocialNetworkOutput, error) {
	socialNetworkAccount, err := u.socialNetworkService.GetSocialNetworkAccount(ctx, input.SocialNetwork)
	if err != nil {
		switch {
		case service.IsValidationError(err):
			return gen.ValidationError{
				Message: err.Error(),
			}, nil
		case service.IsNotFoundError(err):
			return gen.InternalError{
				Message: err.Error(),
			}, nil
		default:
			return nil, ewrap.Errorf(
				"failed to find social socialNetwork %s account: %w", input.SocialNetwork, err,
			)
		}
	}

	pages, err := u.socialNetworkService.GetPagesFromSocialNetwork(socialNetworkAccount)
	if err != nil {
		return gen.InternalError{
			Message: err.Error(),
		}, nil
	}

	var out []*gen.SocialNetworkPage
	for _, page := range pages {
		out = append(out, &gen.SocialNetworkPage{
			PageInfo: &gen.SocialNetworkPageInfo{
				SocialNetworkID: page.ID,
				PageName:        page.Name,
				Description:     &page.Description,
				PreviewImage:    &page.Image,
			},
		})
	}

	return gen.GetPagesFromSocialNetworkResult{
		Pages: out,
	}, nil
}

func (u *SocialNetworkUsecase) CreatePost(
	ctx context.Context,
	input gen.CreatePostInput,
) (gen.CreatePostOutput, error) {
	return nil, nil
}
