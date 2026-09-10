package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/usecase/validation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type partUsecase struct {
	repo   repository.PartRepository
	logger zerolog.Logger
}

func NewPartUsecase(repo repository.PartRepository) PartUsecase {
	return &partUsecase{
		repo:   repo,
		logger: log.With().Str("component", "part_usecase").Logger(),
	}
}

func (u *partUsecase) GetPartByID(ctx context.Context, id int) (*entity.PartRes, error) {
	u.logger.Debug().Int("part_id", id).Msg("GetPartByID")
	return u.repo.GetPartByID(ctx, id)
}

func (u *partUsecase) GetAllParts(ctx context.Context) ([]*entity.PartRes, error) {
	u.logger.Debug().Msg("GetAllParts")
	return u.repo.GetAllParts(ctx)
}

func (u *partUsecase) CreatePart(ctx context.Context, part *entity.PartReq, userID *int) (int64, error) {
	if err := validation.ValidatePartRequest(part, false); err != nil {
		u.logger.Warn().Err(err).Msg("part create validation failed")
		return 0, err
	}

	partID, err := u.repo.CreatePart(ctx, part, userID)
	if err != nil {
		return 0, err
	}

	u.logger.Info().Int64("part_id", partID).Msg("part created successfully")
	return partID, nil
}

func (u *partUsecase) UpdatePart(ctx context.Context, id int, partReq *entity.PartReq, userID int) error {
	if err := validation.ValidatePartRequest(partReq, true); err != nil {
		u.logger.Warn().Err(err).Msg("part update validation failed")
		return err
	}

	if _, err := u.repo.GetPartByID(ctx, id); err != nil {
		return err
	}

	if err := u.repo.UpdatePart(ctx, id, partReq, userID); err != nil {
		return err
	}

	u.logger.Info().Int("part_id", id).Msg("part updated successfully")
	return nil
}

func (u *partUsecase) DeletePart(ctx context.Context, id int, userID int) error {
	if err := u.repo.DeletePart(ctx, id, userID); err != nil {
		return err
	}

	u.logger.Info().Int("part_id", id).Msg("part deleted successfully")
	return nil
}
