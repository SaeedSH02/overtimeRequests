package service

import (
	"context"
	pg "shiftdony/database"
	"shiftdony/models"
	"time"
)

type OvertimeService struct {
	overtimeRepo pg.OvertimeRepository
}

func NewOvertimeService(overtimeRepo pg.OvertimeRepository) *OvertimeService {
	return &OvertimeService{overtimeRepo: overtimeRepo}
}

func (s *OvertimeService) CreateRequest(ctx context.Context, userID, slotID int64) (*models.OvertimeRequest, error) {
	//Check if slot is open
	slot, err := s.overtimeRepo.GetOvertimeSlotByID(ctx, slotID)
	if err != nil {
		return nil, ErrSlotNotFound
	}
	//Check for duplicate
	exists, err := s.overtimeRepo.UserHasPendingRequestForSlot(ctx, userID, slotID)
	if err != nil {
		return nil, ErrRequestNotFound
	}
	if exists {
		return nil, ErrAlreadyApplied
	}

	//Check capacity
	approvedCount, err := s.overtimeRepo.CountApprovedRequestsForSlot(ctx, slotID)
	if err != nil {
		return nil, ErrInternalServer
	}
	if approvedCount >= int(slot.Capacity) {
		slot.Status = "full"
		s.overtimeRepo.UpdateOvertimeSlot(ctx, slot)
		return nil, ErrSlotIsFull
	}
	//new Req
	newRequest := &models.OvertimeRequest{
		UserID:      userID,
		SlotID:      slotID,
		Status:      "pending",
		RequestTime: time.Now(),
	}
	if err := s.overtimeRepo.CreateOvertimeRequest(ctx, newRequest); err != nil {
		return nil, ErrInternalServer
	}

	return newRequest, nil
}

func (s *OvertimeService) GetMyRequests(ctx context.Context, userID int64) ([]models.OvertimeRequest, error) {
	return s.overtimeRepo.GetMyOvertimeRequests(ctx, userID)
}

func (s *OvertimeService) GetAllRequests(ctx context.Context) ([]models.OvertimeRequest, error) {
	return s.overtimeRepo.GetAllOvertimeRequests(ctx)
}

func (s *OvertimeService) GetAvailableSlots(ctx context.Context) ([]models.OvertimeSlot, error) {
	return s.overtimeRepo.GetAvailableOvertimeSlots(ctx)
}

func (s *OvertimeService) UpdateRequestStatus(ctx context.Context, requestID, managerID int64, status string) error {
	request, err := s.overtimeRepo.GetOvertimeRequestByID(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}

	request.Status = status
	request.ReviewedBy = &managerID
	if err := s.overtimeRepo.UpdateOvertimeRequest(ctx, request); err != nil {
		return ErrInternalServer
	}

	if status == "approved" {
		slot, err := s.overtimeRepo.GetOvertimeSlotByID(ctx, request.SlotID)
		if err != nil {
			return ErrInternalServer
		}

		approvedCount, err := s.overtimeRepo.CountApprovedRequestsForSlot(ctx, request.SlotID)
		if err != nil {
			return ErrInternalServer
		}

		if approvedCount >= int(slot.Capacity) {
			slot.Status = "full"
			if err := s.overtimeRepo.UpdateOvertimeSlot(ctx, slot); err != nil {
				return ErrInternalServer
			}
		}
	}

	return nil
}

func (s *OvertimeService) CreateSlot(ctx context.Context, title string, startTime, endTime time.Time, capacity int, creatorID int64) (*models.OvertimeSlot, error) {
	newSlot := &models.OvertimeSlot{
		Title:     title,
		StartTime: startTime,
		EndTime:   endTime,
		Capacity:  int64(capacity),
		CreatedBy: creatorID,
		Status:    "open",
	}

	err := s.overtimeRepo.CreateOvertimeSlot(ctx, newSlot)
	if err != nil {
		return nil, ErrInternalServer
	}

	return newSlot, nil
}

func (s *OvertimeService) GetApprovedRequests(ctx context.Context) ([]models.OvertimeRequest, error) {
	requests, err := s.overtimeRepo.GetApprovedRequests(ctx)
	if err != nil {
		return nil, ErrInternalServer
	}

	for i := range requests {
		if requests[i].User != nil {
			requests[i].User.PasswordHash = ""
		}
	}

	return requests, nil
}

func (s *OvertimeService) GetOvertimeSlots(ctx context.Context) ([]models.OvertimeSlot, error) {
	slots, err := s.overtimeRepo.GetOvertimeSlots(ctx)
	if err != nil {
		return nil, ErrInternalServer
	}
	return slots, nil
}
