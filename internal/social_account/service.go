package social_account

import (
	"amplifr/internal/clients/social_network_client"
	"context"
	"fmt"
)

type SocialAccountService struct {
	repository           SocialAccountRepository
	socialNetworkClients map[string]social_network_client.SocialNetworkClient
}

func NewService(
	repository SocialAccountRepository,
	socialNetworkClients map[string]social_network_client.SocialNetworkClient,
) *SocialAccountService {
	return &SocialAccountService{
		repository:           repository,
		socialNetworkClients: socialNetworkClients,
	}
}

func (sas *SocialAccountService) CreateAccount(network string, credentials string) error {
	var socialAccount SocialAccount
	socialAccount, err := sas.repository.FindAccountByNetwork(context.TODO(), network)
	if err != nil {
		return fmt.Errorf(
			"cannot check account existence: %s",
			err,
		)
	}
	if socialAccount.ID != 0 {
		return fmt.Errorf("account for network %s already exist", network)
	}

	socialAccount = SocialAccount{
		Network:    network,
		Credential: credentials,
	}

	if err := sas.repository.CreateAccount(context.TODO(), &socialAccount); err != nil {
		return fmt.Errorf(
			"cannot create new social account network=%s credentials=%s: %s",
			network,
			credentials,
			err,
		)
	}
	return nil
}

func (sas *SocialAccountService) CreatePost(network, project, post string) error {
	socialAccount, err := sas.repository.FindAccountByNetwork(context.TODO(), network)
	if err != nil {
		return fmt.Errorf(
			"cannot get account: %s",
			err,
		)
	}
	group, err := sas.repository.GetGroup(context.TODO(), socialAccount.ID, project)
	if err != nil {
		return fmt.Errorf(
			"cannot get group for account ID=%d project %s: %s",
			socialAccount.ID,
			project,
			err,
		)
	}
	postId, err := sas.socialNetworkClients[network].CreatePost(
		socialAccount.Credential,
		group.GroupInfo.ID,
		post,
	)
	if err != nil {
		return fmt.Errorf(
			"cannot create new post:\nnetwork %s\nproject %s\n%s",
			network,
			project,
			err,
		)
	}

	fmt.Println(postId)

	return nil
}
