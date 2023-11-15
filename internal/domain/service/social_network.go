package service

import (
	"autoposting/internal/domain/model"
	"autoposting/internal/domain/repository"
	"autoposting/internal/infrastructure/social_network_client"
	"autoposting/internal/presentation/graphql/gen"
	ewrap "autoposting/pkg/err-wrapper"
	"context"
	"log/slog"
	"time"
)

type SocialNetworkService struct {
	logger                          *slog.Logger
	socialNetworkAccountsRepository repository.SocialNetworkAccountsRepository
	socialNetworkPagesRepository    repository.SocialNetworkPagesRepository
	socialNetworkClients            map[model.SocialNetworkName]social_network_client.SocialNetworkClient
}

func NewService(
	logger *slog.Logger,
	socialNetworkAccountsRepository repository.SocialNetworkAccountsRepository,
	socialNetworkPagesRepository repository.SocialNetworkPagesRepository,
	socialNetworkClients map[model.SocialNetworkName]social_network_client.SocialNetworkClient,
) *SocialNetworkService {
	return &SocialNetworkService{
		logger:                          logger,
		socialNetworkAccountsRepository: socialNetworkAccountsRepository,
		socialNetworkPagesRepository:    socialNetworkPagesRepository,
		socialNetworkClients:            socialNetworkClients,
	}
}

func (sns *SocialNetworkService) CreateSocialNetworkAccount(
	ctx context.Context,
	socialNetwork string,
	credential string,
) error {
	socialNetworkName, err := getSocialNetworkName(socialNetwork)
	if err != nil {
		sns.logger.Error(err.Error())
		return err
	}

	socialNetworkAccount := &model.SocialNetworkAccount{
		SocialNetwork: socialNetworkName,
		Credentials:   credential,
	}

	return sns.socialNetworkAccountsRepository.CreateAccount(ctx, socialNetworkAccount)
}

func (sns *SocialNetworkService) GetSocialNetworkAccount(
	ctx context.Context,
	socialNetwork string,
) (*model.SocialNetworkAccount, error) {
	socialNetworkName, err := getSocialNetworkName(socialNetwork)
	if err != nil {
		return nil, err
	}

	return sns.socialNetworkAccountsRepository.FindBySocialNetwork(ctx, socialNetworkName)
}

func (sns *SocialNetworkService) CreateSocialNetworkPage(
	ctx context.Context,
	input gen.CreateSocialNetworkPageInput,
) error {
	socialNetworkPage := &model.SocialNetworkPage{
		AccountID: input.SocialNetworkAccountID,
		Project:   input.Project,
		PageID:    input.PageInfo.SocialNetworkID,
		PageInfo: &model.SocialNetworkPageInfo{
			Name:         input.PageInfo.PageName,
			Description:  *input.PageInfo.Description,
			PreviewImage: *input.PageInfo.PreviewImage,
		},
	}
	if input.AccessToken != nil {
		socialNetworkPage.AccessToken = &model.AccessToken{
			Token:     input.AccessToken.Token,
			ExpiresIn: *input.AccessToken.ExpiresIn,
		}
	}

	return sns.socialNetworkPagesRepository.CreatePage(ctx, socialNetworkPage)
}

func (sns *SocialNetworkService) GetAuthURL(
	socialNetworkName model.SocialNetworkName,
	credentials string,
) (string, error) {
	authUrl, err := sns.socialNetworkClients[socialNetworkName].GetAuthURL(credentials)
	if err != nil {
		return "", NewInternalError(err.Error())
	}

	return authUrl, nil
}

func (sns *SocialNetworkService) GetAccessToken(
	ctx context.Context,
	socialNetworkAccount *model.SocialNetworkAccount,
	params map[string][]string,
) (string, error) {
	token, err := sns.socialNetworkClients[socialNetworkAccount.SocialNetwork].GetAccessToken(
		socialNetworkAccount.Credentials,
		params,
	)
	if err != nil {
		return "", ewrap.Errorf(
			"failed to get access token from social network %s: %w",
			socialNetworkAccount.SocialNetwork,
			err,
		)
	}

	return token, nil
}

func (sns *SocialNetworkService) SaveAccessToken(
	ctx context.Context,
	token string,
	socialNetworkAccount *model.SocialNetworkAccount,
) error {
	socialNetworkAccount.AccessToken = &model.AccessToken{
		Token:     token,
		ExpiresIn: getTokenExpires(&socialNetworkAccount.SocialNetwork),
	}

	if _, err := sns.socialNetworkAccountsRepository.UpdateAccount(ctx, socialNetworkAccount); err != nil {
		return ewrap.Errorf(
			"failed to set token for social network %s account: %w",
			socialNetworkAccount.SocialNetwork,
			err,
		)
	}

	return nil
}

func (sns *SocialNetworkService) GetPagesFromSocialNetwork(
	socialNetworkAccount *model.SocialNetworkAccount,
) ([]social_network_client.SocialNetworkPage, error) {
	pages, err := sns.socialNetworkClients[socialNetworkAccount.SocialNetwork].GetAccountPages(
		socialNetworkAccount.Credentials,
		socialNetworkAccount.AccessToken.Token,
	)

	if err != nil {
		return nil, ewrap.Errorf(
			"failed to get pages from social network %s: %w",
			socialNetworkAccount.SocialNetwork,
			err,
		)
	}

	return pages, nil
}

func (sns *SocialNetworkService) CreatePost(network, project, post string) error {
	return nil
}

func getSocialNetworkName(socialNetwork string) (model.SocialNetworkName, error) {
	socialNetworkName := model.SocialNetworkName(socialNetwork)
	if err := socialNetworkName.Validate(); err != nil {
		return "", NewValidationError(
			err.Error(),
			"socialNetworkName",
			"empty",
		)
	}
	return socialNetworkName, nil
}

func getTokenExpires(socialNetworkName *model.SocialNetworkName) string {
	switch *socialNetworkName {
	case model.FB:
		// У FACEBOOK longLiveToken живет 60 дней
		return time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339)
	default:
		return ""
	}
}
