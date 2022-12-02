package social_account

import (
	"autoposting/internal/clients/social_network_client"
	logging "autoposting/pkg"
	"context"
	"fmt"
)

type SocialAccountService struct {
	logger               *logging.Logger
	repository           SocialAccountRepository
	socialNetworkClients map[string]social_network_client.SocialNetworkClient
}

func NewService(
	logger *logging.Logger,
	repository SocialAccountRepository,
	socialNetworkClients map[string]social_network_client.SocialNetworkClient,
) *SocialAccountService {
	return &SocialAccountService{
		logger:               logger,
		repository:           repository,
		socialNetworkClients: socialNetworkClients,
	}
}

func (sas *SocialAccountService) GetAuthURL(network string) (string, error) {
	var socialAccount SocialAccount
	socialAccount, err := sas.repository.FindAccountByNetwork(context.TODO(), network)
	if err != nil {
		sas.logger.Errorf("cannot get social account by network %s:\n%s", network, err)
		return "", err
	}

	authURL, err := sas.socialNetworkClients[network].GetAuthURL(socialAccount.Credential)
	if err != nil {
		sas.logger.Errorf("cannot get auth URL for %s:\n%s", network, err)
		return "", err
	}

	return authURL, nil
}

func (sas *SocialAccountService) GetAccessToken(queryParams map[string][]string) (string, error) {
	var socialAccount SocialAccount
	network := queryParams["network"][0]
	delete(queryParams, "network")

	socialAccount, err := sas.repository.FindAccountByNetwork(context.TODO(), network)
	if err != nil {
		sas.logger.Errorf("cannot get social account by network %s:\n%s", network, err)
		return "", err
	}

	accessToken, err := sas.socialNetworkClients[network].GetAccessToken(socialAccount.Credential, queryParams)
	if err != nil {
		sas.logger.Errorf("cannot get auth URL for %s:\n%s", network, err)
		return "", err
	}

	return accessToken, nil
}

func (sas *SocialAccountService) CreateAccount(network string, credentials string) error {
	var socialAccount SocialAccount
	socialAccount, err := sas.repository.FindAccountByNetwork(context.TODO(), network)
	if err != nil {
		sas.logger.Errorf("cannot get social account by network %s:\n%s", network, err)
		return err
	}
	if socialAccount.ID != 0 {
		sas.logger.Errorf("account for network %s already exist", network)
		return err
	}

	socialAccount = SocialAccount{
		Network:    network,
		Credential: credentials,
	}

	if err := sas.repository.CreateAccount(context.TODO(), &socialAccount); err != nil {
		sas.logger.Errorf(
			"cannot create new social account network=%s credentials=%s: %s",
			network,
			credentials,
			err,
		)
		return err
	}
	return nil
}

func (sas *SocialAccountService) CreatePost(network, project, post string) error {
	socialAccount, err := sas.repository.FindAccountByNetwork(context.TODO(), network)
	if err != nil {
		sas.logger.Errorf("cannot get account: %s", err)
		return err
	}
	group, err := sas.repository.GetGroup(context.TODO(), socialAccount.ID, project)
	if err != nil {
		sas.logger.Errorf(
			"cannot get group for account ID=%d project %s: %s",
			socialAccount.ID,
			project,
			err,
		)
		return err
	}
	postId, err := sas.socialNetworkClients[network].CreatePost(
		socialAccount.Credential,
		group.GroupInfo.ID,
		post,
	)
	if err != nil {
		sas.logger.Errorf(
			"cannot create new post:\nnetwork %s\nproject %s\n%s",
			network,
			project,
			err,
		)
		return err
	}

	fmt.Println(postId)

	return nil
}
