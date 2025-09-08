package repository

import (
	"context"
	"shiftdony/models"

	"github.com/uptrace/bun"
)

type overtimeRepository struct {
	db *bun.DB
}

func NewOvertimeRepository(db *bun.DB) *overtimeRepository {
	return &overtimeRepository{db: db}
}

func (r *overtimeRepository) CreateOvertimeSlot(ctx context.Context, slot *models.OvertimeSlot) error {
	_, err := r.db.NewInsert().Model(slot).Exec(ctx)
	return err
}

func (r *overtimeRepository) GetOvertimeSlots(ctx context.Context) ([]models.OvertimeSlot, error) {
	var slots []models.OvertimeSlot
	err := r.db.NewSelect().
		Model(&slots).
		Relation("Creator", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column("full_name")
		}).
		Order("start_time DESC").
		Scan(ctx)
	return slots, err
}

func (r *overtimeRepository) GetAvailableOvertimeSlots(ctx context.Context) ([]models.OvertimeSlot, error) {
	var slots []models.OvertimeSlot
	err := r.db.NewSelect().
		Model(&slots).
		Relation("Creator", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("full_name")
		}).
		Where("status = ?", "open").
		Order("start_time ASC").
		Scan(ctx)

	return slots, err
}

func (r *overtimeRepository) GetOvertimeSlotByID(ctx context.Context, slotID int64) (*models.OvertimeSlot, error) {
	var slot models.OvertimeSlot
	err := r.db.NewSelect().
		Model(&slot).
		Where("id = ? AND status = ?", slotID, "open").
		Scan(ctx)
	return &slot, err
}

func (r *overtimeRepository) UserHasPendingRequestForSlot(ctx context.Context, userID, slotID int64) (bool, error) {
	return r.db.NewSelect().
		Model((*models.OvertimeRequest)(nil)).
		Where("user_id = ? AND slot_id = ?", userID, slotID).
		Exists(ctx)
}

func (r *overtimeRepository) CountApprovedRequestsForSlot(ctx context.Context, slotID int64) (int, error) {
	return r.db.NewSelect().
		Model((*models.OvertimeRequest)(nil)).
		Where("slot_id = ? AND status = ?", slotID, "approved").
		Count(ctx)
}

func (r *overtimeRepository) CreateOvertimeRequest(ctx context.Context, req *models.OvertimeRequest) error {
	_, err := r.db.NewInsert().Model(req).Exec(ctx)
	return err
}

func (r *overtimeRepository) UpdateOvertimeSlot(ctx context.Context, slot *models.OvertimeSlot) error {
	_, err := r.db.NewUpdate().Model(slot).WherePK().Exec(ctx)
	return err
}

func (r *overtimeRepository) GetMyOvertimeRequests(ctx context.Context, userID int64) ([]models.OvertimeRequest, error) {
	var requests []models.OvertimeRequest
	err := r.db.NewSelect().
		Model(&requests).
		Where("user_id = ?", userID).
		Relation("Slot").
		Order("request_time DESC").
		Scan(ctx)
	return requests, err
}

func (r *overtimeRepository) GetAllOvertimeRequests(ctx context.Context) ([]models.OvertimeRequest, error) {
	var requsets []models.OvertimeRequest
	err := r.db.NewSelect().
		Model(&requsets).
		Relation("User").
		Relation("Slot").
		Order("request_time DESC").
		Scan(ctx)
	return requsets, err
}

func (r *overtimeRepository) GetOvertimeRequestByID(ctx context.Context, requestID int64) (*models.OvertimeRequest, error) {
	var request models.OvertimeRequest
	err := r.db.NewSelect().
		Model(&request).
		Where("id = ?", requestID).
		For("UPDATE").
		Scan(ctx)
	return &request, err
}

func (r *overtimeRepository) UpdateOvertimeRequest(ctx context.Context, req *models.OvertimeRequest) error {
	_, err := r.db.NewUpdate().Model(req).WherePK().Exec(ctx)
	return err
}

func (r *overtimeRepository) GetApprovedRequests(ctx context.Context) ([]models.OvertimeRequest, error) {
	var approvedRequests []models.OvertimeRequest
	err := r.db.NewSelect().
		Model(&approvedRequests).
		Relation("User").Relation("Slot").
		Where("?TableAlias.status = ?", "approved").
		Order("slot.start_time ASC").Scan(ctx)
	return approvedRequests, err
}
