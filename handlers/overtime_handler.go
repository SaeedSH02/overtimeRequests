package handlers

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	postgres "shiftdony/database"
	log "shiftdony/logs"
	"shiftdony/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type OvertimeHandler struct {
	db *postgres.Postgres
}


func NewOvertimeHandler(db *postgres.Postgres) *OvertimeHandler {
	return &OvertimeHandler{db: db}
}

func (h *OvertimeHandler) CreateOvertimeSlot(c *gin.Context) {
	var input CreateOvertimeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid Input",
			"reason":  "INVALID_INPUT"})
		return
	}

	creatorID, _ := c.Get("userID")
	newSlot := models.OvertimeSlot{
		Title:     input.Title,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Capacity:  int64(input.Capacity),
		CreatedBy: int64(creatorID.(float64)),
		Status:    "open",
	}
	_, err := h.db.DB().
		NewInsert().
		Model(&newSlot).
		Exec(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create overtime slot",
			"reason":  "SERVER_ERROR"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "new slot": newSlot})

}

func (h *OvertimeHandler) GetOvertimeSlots(c *gin.Context) {
	var slots []models.OvertimeSlot

	err := h.db.DB().
		NewSelect().
		Model(&slots).
		Relation("Creator", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column("full_name")
		}).
		Order("start_time DESC").
		Scan(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch overtime slots",
			"reason":  "SERVER_ERROR"})
		return
	}
	if slots == nil {
		slots = make([]models.OvertimeSlot, 0)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "slots": slots})
}

//Get Available OvertimeSlots

func (h *OvertimeHandler) GetAvailableOvertimeSlots(c *gin.Context) {
	var slots []models.OvertimeSlot

	err := h.db.DB().NewSelect().
		Model(&slots).
		Where("status = ?", "open").
		Order("start_time ASC").
		Scan(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch available slots",
			"reason":  "SERVER_ERROR"})
		log.Gl.Info("Failed to fetch available slots" + err.Error())
		return
	}
	if slots == nil {
		slots = make([]models.OvertimeSlot, 0)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "slots": slots})
}

func (h *OvertimeHandler) CreateOvertimeRequest(c *gin.Context) {
	var input CreateRequestInput
	// ctx := context.Background()
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid Input",
			"reason":  "INVALID_INPUT"})
		return
	}
	userID, _ := c.Get("userID")
	userID64 := int64(userID.(float64))

	//check Opens shift
	var slot models.OvertimeSlot
	err := h.db.DB().
		NewSelect().Model(&slot).
		Where("id = ? AND status = ?", input.SlotID, "open").
		Scan(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Slot not found or is not open for requests",
			"reason":  "ALREADY_EXISTS"})
		return
	}

	// Verify user duplicate shift request
	exists, err := h.db.DB().NewSelect().
		Model((*models.OvertimeRequest)(nil)).
		Where("user_id = ? AND slot_id = ?", userID64, input.SlotID).
		Exists(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch overtime slots",
			"reason":  "SERVER_ERROR"})
		log.Gl.Info("Database error while checking existing request" + err.Error())
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "You have already applied for this slot",
			"reason":  "ALREADY_EXISTS"})
		return
	}

	// Get approved requests count for this overTime
	approvedCount, err := h.db.DB().NewSelect().
		Model((*models.OvertimeRequest)(nil)).
		Where("slot_id = ? AND status = ?", input.SlotID, "approved").
		Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
			"reason":  "SERVER_ERROR"})
		log.Gl.Info("Database Error while checking capacity" + err.Error())
		return
	}
	if approvedCount >= int(slot.Capacity) {
		slot.Status = "full"
		h.db.DB().NewUpdate().Model(&slot).WherePK().Exec(c.Request.Context())
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "This overtime slot is already full",
			"reason":  "ALREADY_EXISTS"})
		return
	}

	//Make New Req
	newRequest := models.OvertimeRequest{
		UserID:      userID64,
		SlotID:      input.SlotID,
		Status:      "pending",
		RequestTime: time.Now(),
	}
	_, err = h.db.DB().NewInsert().Model(&newRequest).Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to creat request",
			"reason":  "SERVER_ERROR"})
		log.Gl.Info("Failed to create request" + err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Request submitted successfully", "request": newRequest})
}

// Get My ocertime requests
func (h *OvertimeHandler) GetMyOvertimeRequests(c *gin.Context) {
	userID, _ := c.Get("userID")
	userID64 := int64(userID.(float64))

	var requests []models.OvertimeRequest

	err := h.db.DB().NewSelect().Model(&requests).
		Where("user_id = ?", userID64).
		Relation("Slot").
		Order("request_time DESC").
		Scan(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
			"reason":  "SERVER_ERROR"})
		log.Gl.Info("Failed to fetch user requests", zap.Int64("user id:", userID64))
		return
	}

	if requests == nil {
		requests = make([]models.OvertimeRequest, 0)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "my_request": requests})
}

// Get All overtime Req for Admins
func (h *OvertimeHandler) GetAllOvertimeRequests(c *gin.Context) {
	var requests []models.OvertimeRequest

	err := h.db.DB().NewSelect().Model(&requests).
		Relation("User").
		Relation("Slot").
		Order("request_time DESC").
		Scan(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
			"reason":  "SERVER_ERROR"})
		log.Gl.Error("Failed to fetch requests for admin: ", zap.Error(err))
		return
	}
	for i := range requests {
		if requests[i].User != nil {
			requests[i].User.PasswordHash = ""

		}
	}
	if requests == nil {
		requests = make([]models.OvertimeRequest, 0)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "all_requests": requests})
}

// Update Request Status
func (h *OvertimeHandler) UpdateOvertimeReqStatus(c *gin.Context) {
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Bad Request",
			"reason":  "INVALID_INPUT"})
		log.Gl.Info("bad request detected: ", zap.Error(err))
		return
	}
	var input UpdateRequestStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Bad Request",
			"reason":  "INVALID_INPUT"})
		log.Gl.Info("Bad request detected: ", zap.Error(err))
		return
	}

	managerID, exists := c.Get("userID")
	if !exists || managerID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Internal Server Error",
			"reason":  "SERVER_ERROR"})
		log.Gl.Error("managerID not found", zap.Error(err))
		return
	}
	managerID64 := int64(managerID.(float64))

	tx, err := h.db.DB().BeginTx(c, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
			"reason":  "SERVER_ERROR"})
		log.Gl.Error("Failed to start transaction", zap.Int64("Admin Id", managerID64), zap.Error(err))
		return
	}
	defer tx.Rollback() //If an error occurs, all changes will be rolled back

	var request models.OvertimeRequest
	err = tx.NewSelect().Model(&request).
		Where("id = ?", requestID).
		For("UPDATE").
		Scan(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Request not found",
			"reason":  "NOT_FOUND"})
		log.Gl.Info("Request not found ", zap.Error(err))
		return
	}
	request.Status = input.Status
	request.ReviewedBy = &managerID64
	_, err = tx.NewUpdate().Model(&request).WherePK().Exec(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
			"reason":  "SERVER_ERROR"})
		log.Gl.Error("Failed to update request status", zap.Int64("Admin Id", managerID64), zap.Error(err))
		return
	}

	if input.Status == "approved" {
		var slot models.OvertimeSlot
		tx.NewSelect().Model(&slot).Where("id = ?", request.SlotID).Scan(c)

		approvedCount, _ := tx.NewSelect().Model((*models.OvertimeRequest)(nil)).
			Where("slot_id = ? AND status = ?", request.SlotID, "approved").
			Count(c)

		if approvedCount >= int(slot.Capacity) {
			//if capacity is full change status to Full
			_, err = tx.NewUpdate().Model((*models.OvertimeSlot)(nil)).
				Set("status = ?", "full").
				Where("id = ?", request.SlotID).
				Exec(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "Internal Server Error",
					"reason":  "SERVER_ERROR"})
				log.Gl.Error("Failed to update slot status to full", zap.Int64("Admin Id", managerID64), zap.Error(err))
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
			"reason":  "SERVER_ERROR"})
		log.Gl.Error("Failed to commit transaction", zap.Int64("Admin Id", managerID64), zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Request status updated successfully"})
}

// Export Approved Requests As CSV
func (h *OvertimeHandler) ExportApprovedRequestsAsCSV(c *gin.Context) {
	var approvedRequests []models.OvertimeRequest

	// Relation("Creator", func(sq *bun.SelectQuery) *bun.SelectQuery {
	// 	return sq.Column("full_name")
	// }).
	err := h.db.DB().NewSelect().
		Model(&approvedRequests).
		Relation("User").Relation("Slot").
		Where("?TableAlias.status = ?", "approved").
		Order("slot.start_time ASC").
		Scan(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch approved requests",
			"reason":  "SERVER_ERROR"})
		log.Gl.Error("Failed to fetch approved requests", zap.Error(err))
		return
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="overtime_report.csv"`)

	writer := csv.NewWriter(c.Writer)

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
