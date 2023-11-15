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
	"github.com/ztrue/tracerr"
)

type SocialNetworkPagesRepository struct {
	db dao.Querier
}

func NewSocialNetworkPagesRepository(db dao.Querier) *SocialNetworkPagesRepository {
	return &SocialNetworkPagesRepository{db: db}
}

func (s SocialNetworkPagesRepository) CreatePage(
	ctx context.Context,
	socialNetworkPage *model.SocialNetworkPage,
) error {
	isExist, err := s.db.NewSelect().
		Model(socialNetworkPage).
		Where(`"account_id" = ?`, socialNetworkPage.AccountID).
		Where(`"page_id" = ?`, socialNetworkPage.PageID).
		Exists(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return ewrap.Errorf("failed to check page exists: %w", err)
	}

	_, err = s.db.NewInsert().
		Model(socialNetworkPage).
		On(`CONFLICT ON CONSTRAINT "SOCIAL_NETWORK_PAGES_UNIQUE" DO UPDATE`).
		Returning("id").
		Exec(ctx)
	if err != nil {
		return tracerr.Errorf("failed to create social network page: %w", err)
	}
	if isExist {
		return service.NewPageAlreadyExistsError(
			fmt.Sprintf("social network page already exist, page with id=%d updated", socialNetworkPage.ID),
		)
	}

	return nil
}

func (s SocialNetworkPagesRepository) FindPage() {
	//TODO implement me
	panic("implement me")
}
