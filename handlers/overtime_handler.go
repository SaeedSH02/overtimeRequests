package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	log "shiftdony/logs"
	"shiftdony/models"
	"shiftdony/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OvertimeHandler struct {
	overtimeService *service.OvertimeService
}

func NewOvertimeHandler(overtimeService *service.OvertimeService) *OvertimeHandler {
	return &OvertimeHandler{overtimeService: overtimeService}
}

func (h *OvertimeHandler) CreateOvertimeSlot(c *gin.Context) {
	var input CreateOvertimeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid input data", "INVALID_INPUT")
		return
	}

	creatorIDVal, _ := c.Get("userID")
	creatorID := int64(creatorIDVal.(float64))

	newSlot, err := h.overtimeService.CreateSlot(
		c.Request.Context(),
		input.Title,
		input.StartTime,
		input.EndTime,
		input.Capacity,
		creatorID,
	)

	if err != nil {
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to create overtime slot", "SERVER_ERROR")
		return
	}

	SendSuccessResponse(c, http.StatusCreated, newSlot)
}

func (h *OvertimeHandler) GetOvertimeSlots(c *gin.Context) {
	slots, err := h.overtimeService.GetOvertimeSlots(c.Request.Context())

	if err != nil {
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch overtime slots", "SERVER_ERROR")
		return
	}

	if slots == nil {
		slots = make([]models.OvertimeSlot, 0)
	}

	responseSlots := make([]SlotResponse, 0, len(slots))
	for _, slot := range slots {
		var creatorName string
		if slot.Creator != nil {
			creatorName = slot.Creator.FullName
		}

		responseSlots = append(responseSlots, SlotResponse{
			ID:        slot.ID,
			Title:     slot.Title,
			StartTime: slot.StartTime,
			EndTime:   slot.EndTime,
			Capacity:  slot.Capacity,
			Status:    slot.Status,
			Creator:   creatorName,
		})
	}

	SendSuccessResponse(c, http.StatusOK, responseSlots)
}

// Get Available OvertimeSlots
func (h *OvertimeHandler) GetAvailableOvertimeSlots(c *gin.Context) {

	slots, err := h.overtimeService.GetAvailableSlots(c.Request.Context())

	if err != nil {
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch available slots", "SERVER_ERROR")
		log.Gl.Error("Failed to fetch available slots", zap.Error(err))
		return
	}
	if slots == nil {
		slots = make([]models.OvertimeSlot, 0)
	}

	SendSuccessResponse(c, http.StatusOK, slots)

}

func (h *OvertimeHandler) CreateOvertimeRequest(c *gin.Context) {
	var input CreateRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid inpud data", "INVALID_INPUT")
		return
	}
	userIDVal, _ := c.Get("userID")
	userID := int64(userIDVal.(float64))

	newRequest, err := h.overtimeService.CreateRequest(c.Request.Context(), userID, input.SlotID)

	if err != nil {
		switch err {
		case service.ErrSlotNotFound:
			SendErrorResponse(c, http.StatusNotFound, "Slot not found or is not open for requests", "SLOT_NOT_FOUND")
		case service.ErrAlreadyApplied:
			SendErrorResponse(c, http.StatusConflict, "You have already applied for this slot", "ALREADY_APPLIED")
		case service.ErrSlotIsFull:
			SendErrorResponse(c, http.StatusConflict, "This overtime slot is already full", "SLOT_FULL")
		default:
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to create request", "SERVER_ERROR")
			log.Gl.Error("Failed to create request", zap.Error(err))
		}
		return
	}

	SendSuccessResponse(c, http.StatusOK, newRequest)
}

// Get My ocertime requests
func (h *OvertimeHandler) GetMyOvertimeRequests(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := int64(userIDVal.(float64))

	requests, err := h.overtimeService.GetMyRequests(c.Request.Context(), userID)

	if err != nil {
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch your requests", "SERVER_ERROR")
		log.Gl.Error("Failed to fetch user requests", zap.Error(err))
		return
	}

	if requests == nil {
		requests = make([]models.OvertimeRequest, 0)
	}

	SendSuccessResponse(c, http.StatusOK, requests)
}

// Get All overtime Req for Admins
func (h *OvertimeHandler) GetAllOvertimeRequests(c *gin.Context) {

	requests, err := h.overtimeService.GetAllRequests(c.Request.Context())

	if err != nil {
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch all requests", "SERVER_ERROR")
		log.Gl.Error("Failed to fetch all requests", zap.Error(err))
		return
	}

	if requests == nil {
		requests = make([]models.OvertimeRequest, 0)
	}

	SendSuccessResponse(c, http.StatusOK, requests)
}

// Update Request Status
func (h *OvertimeHandler) UpdateOvertimeReqStatus(c *gin.Context) {

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid request ID format", "INVALID_INPUT")
		return
	}
	var input UpdateRequestStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid input data for status", "INVALID_INPUT")
		return
	}

	managerIDVal, _ := c.Get("userID")
	managerID := int64(managerIDVal.(float64))

	err = h.overtimeService.UpdateRequestStatus(c.Request.Context(), requestID, managerID, input.Status)

	if err != nil {
		if err == service.ErrRequestNotFound {
			SendErrorResponse(c, http.StatusNotFound, "The requested overtime request was not found", "NOT_FOUND")
		} else {
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to update request status", "SERVER_ERROR")
		}
		return
	}

	SendSuccessResponse(c, http.StatusOK, gin.H{
		"message": "Request status updated successfully",
	})
}

// Export Approved Requests As CSV
func (h *OvertimeHandler) ExportApprovedRequestsAsCSV(c *gin.Context) {
	approvedRequests, err := h.overtimeService.GetApprovedRequests(c.Request.Context())
	if err != nil {
		SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch approved requests", "SERVER_ERROR")
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="approved_requests_report.csv"`)

	writer := csv.NewWriter(c.Writer)

	header := []string{"RequestID", "PersonnelCode", "FullName", "TeamID", "SlotTitle", "StartTime", "EndTime", "ReviewedByManagerID"}
	writer.Write(header)

	for _, req := range approvedRequests {
		var managerIDStr string
		if req.ReviewedBy != nil {
			managerIDStr = fmt.Sprintf("%d", *req.ReviewedBy)
		}

		row := []string{
			fmt.Sprintf("%d", req.ID),
			req.User.PersonnelCode,
			req.User.FullName,
			fmt.Sprintf("%d", req.User.TeamID),
			req.Slot.Title,
			req.Slot.StartTime.Format(time.RFC3339),
			req.Slot.EndTime.Format(time.RFC3339),
			managerIDStr,
		}
		writer.Write(row)
	}

	writer.Flush()
}
