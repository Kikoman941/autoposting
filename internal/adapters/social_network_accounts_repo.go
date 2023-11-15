package adapters

import (
	"autoposting/internal/domain/model"
	"autoposting/internal/domain/service"
	"autoposting/internal/infrastructure/sqlc-pg/dao"
	ewrap "autoposting/pkg/err-wrapper"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
)

type SocialNetworkAccountsRepository struct {
	db dao.Querier
}

type FindSocialNetworkAccountQuery struct {
	SocialNetworkAnyOf []model.SocialNetworkName
}

func NewSocialNetworkAccountsRepository(db dao.Querier) *SocialNetworkAccountsRepository {
	return &SocialNetworkAccountsRepository{db: db}
}

func (s SocialNetworkAccountsRepository) CreateAccount(
	ctx context.Context,
	socialNetworkAccount *model.SocialNetworkAccount,
) error {
	createSocialNetworkAccountParams := dao.CreateSocialNetworkAccountParams{
		SocialNetwork: string(socialNetworkAccount.SocialNetwork),
		Credentials:   socialNetworkAccount.Credentials,
		AccessToken:   *socialNetworkAccount.AccessToken,
	}

	rows, err := s.db.GetSocialNetworkAccounts(ctx, createSocialNetworkAccountParams)
	if err != nil {
		return nil, wrap.Errorf("failed to find providers: %w", err)
	}
	//isExist, err := s.db.NewSelect().
	//	Model(socialNetworkAccount).
	//	Where(`"social_network" = ?`, socialNetworkAccount.SocialNetwork).
	//	Exists(ctx)
	//if err != nil && !errors.Is(err, sql.ErrNoRows) {
	//	return ewrap.Errorf("failed to check account exists: %w", err)
	//}
	//if isExist {
	//	return service.NewSocialNetworkAccountAlreadyExistsError(
	//		"social network account already exists",
	//	)
	//}
	//
	//_, err = s.db.NewInsert().
	//	Model(socialNetworkAccount).
	//	Returning("id").
	//	Exec(ctx)
	//if err != nil || socialNetworkAccount.ID == 0 {
	//	return ewrap.Errorf("failed to create account: %w", err)
	//}
	//return nil
}

func (s SocialNetworkAccountsRepository) FindAccounts(ctx context.Context, query FindSocialNetworkAccountQuery) ([]model.SocialNetworkAccount, error) {
	var accountRows []model.SocialNetworkAccount
	q := s.db.NewSelect().Model(&accountRows)

	if len(query.SocialNetworkAnyOf) != 0 {
		q.Where("social_network IN (?)", bun.In(query.SocialNetworkAnyOf))
	}

	if err := q.Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return accountRows, nil
		}
		return nil, ewrap.Errorf("failed to select social network account: %s", err)
	}
	return accountRows, nil
}

func (s SocialNetworkAccountsRepository) UpdateAccount(
	ctx context.Context,
	socialNetworkAccount *model.SocialNetworkAccount,
) (*model.SocialNetworkAccount, error) {
	_, err := s.db.NewUpdate().
		Model(socialNetworkAccount).
		WherePK().
		Exec(ctx)
	if err != nil {
		return nil, ewrap.Errorf("failed to update social network account: %s", err)
	}
	return socialNetworkAccount, nil
}

func (s SocialNetworkAccountsRepository) FindBySocialNetwork(
	ctx context.Context,
	socialNetwork model.SocialNetworkName,
) (*model.SocialNetworkAccount, error) {
	accounts, err := s.FindAccounts(ctx, FindSocialNetworkAccountQuery{
		SocialNetworkAnyOf: []model.SocialNetworkName{
			socialNetwork,
		},
	})

	if err != nil {
		if len(accounts) == 0 {
			return nil, service.NewNotFoundError(
				fmt.Sprintf("social network %v accounts not found", socialNetwork),
			)
		}
		return nil, err
	}

	return &accounts[0], nil
}
