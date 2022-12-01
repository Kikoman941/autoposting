package storage

import (
	"autoposting/internal/clients/db"
	"autoposting/internal/social_account"
	"context"
	"errors"
	"github.com/go-pg/pg/v10"
)

type repository struct {
	client db.Client
}

func NewRepository(dbClient db.Client) social_account.SocialAccountRepository {
	return &repository{
		client: dbClient,
	}
}

func (r *repository) CreateAccount(ctx context.Context, sa *social_account.SocialAccount) error {
	query := r.client.ModelContext(ctx, sa)
	_, err := query.Insert()
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetOneAccount(ctx context.Context, id int) (social_account.SocialAccount, error) {
	socialAccount := &social_account.SocialAccount{
		ID: id,
	}
	err := r.client.ModelContext(ctx, &socialAccount).WherePK().Select()
	if err != nil {
		return social_account.SocialAccount{}, err
	}
	return *socialAccount, nil
}

func (r *repository) FindAccountByNetwork(ctx context.Context, network string) (social_account.SocialAccount, error) {
	socialAccount := social_account.SocialAccount{}
	err := r.client.ModelContext(ctx, &socialAccount).
		Where(`"network" = ?`, network).
		Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return socialAccount, nil
		}
		return socialAccount, err
	}
	return socialAccount, nil
}

func (r *repository) GetGroup(ctx context.Context, accountId int, project string) (social_account.Group, error) {
	group := social_account.Group{}
	err := r.client.ModelContext(ctx, &group).
		Where(`"account_id" = ? and "project" = ?`, accountId, project).
		Select()
	if err != nil {
		return social_account.Group{}, err
	}
	return group, nil
}
