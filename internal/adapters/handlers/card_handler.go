package handlers

import (
	"interview-tracker/internal/entities"
	"interview-tracker/internal/middleware"
	"interview-tracker/internal/models/card_models"
	"interview-tracker/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CardHandler struct{ uc *usecases.CardUsecase }

func NewCardHandler(uc *usecases.CardUsecase) *CardHandler { return &CardHandler{uc} }

// @Summary List cards
// @Tags cards
// @Security BearerAuth
// @Produce json
// @Param page       query int    false "page"
// @Param page_size  query int    false "page size"
// @Param status     query string false "Filter status" Enums(todo,in_progress,done) default(todo)
// @Success 200 {object} map[string]any
// @Router /interview-tracker/authen/cards [get]
func (h *CardHandler) List(c *gin.Context) {
	var q card_models.ListCardsQuery
	_ = c.ShouldBindQuery(&q)
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 || q.PageSize > 100 {
		q.PageSize = 10
	}

	items, total, err := h.uc.List(q.Status, q.Page, q.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "page": q.Page, "page_size": q.PageSize})
}

// @Summary Get card detail
// @Tags cards
// @Security BearerAuth
// @Produce json
// @Param id path string true "card id"
// @Success 200 {object} map[string]any
// @Router /interview-tracker/authen/cards/{id} [get]
func (h *CardHandler) Detail(c *gin.Context) {
	id := c.Param("id")
	card, err := h.uc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"card": card})
}

// @Summary Create card
// @Tags cards
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body card_models.CreateCardReq true "card create JSON"
// @Success 200 {object} map[string]any
// @Router /interview-tracker/authen/cards [post]
func (h *CardHandler) Create(c *gin.Context) {
	var req card_models.CreateCardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := middleware.GetSession(c)
	card := &entities.Card{
		Title:         req.Title,
		Description:   req.Description,
		CandidateName: session.Email,
		ScheduledAt:   req.ScheduledAt,
		CreatedBy:     session.UserID,
	}
	if err := h.uc.Create(card); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"card": card})
}

// @Summary Update card (partial)
// @Tags cards
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "card id"
// @Param request body card_models.UpdateCardReq true "patch"
// @Success 200 {object} map[string]any
// @Router /interview-tracker/authen/cards/{id} [patch]
func (h *CardHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req card_models.UpdateCardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	patch := map[string]any{}
	if req.Title != nil {
		patch["title"] = *req.Title
	}
	if req.Description != nil {
		patch["description"] = *req.Description
	}
	if req.ScheduledAt != nil {
		patch["scheduled_at"] = *req.ScheduledAt
	}
	card, err := h.uc.UpdatePartial(id, patch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"card": card})
}

// @Summary Change status
// @Tags cards
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "card id"
// @Param request body card_models.UpdateCardStatusReq true "new status"
// @Success 200 {object} map[string]any
// @Router /interview-tracker/authen/cards/{id}/status [patch]
func (h *CardHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req card_models.UpdateCardStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	session := middleware.GetSession(c)
	card, err := h.uc.UpdateStatus(id, req.Status, session.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"card": card})
}

// @Summary Add comment
// @Tags comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "card id"
// @Param request body card_models.AddCommentReq true "comment"
// @Success 200 {object} map[string]any
// @Router /interview-tracker/authen/cards/{id}/comments [post]
func (h *CardHandler) AddComment(c *gin.Context) {
	id := c.Param("id")
	var req card_models.AddCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	session := middleware.GetSession(c)
	if err := h.uc.AddComment(session.UserID, id, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// @Summary List comments
// @Tags comments
// @Security BearerAuth
// @Produce json
// @Param id path string true "card id"
// @Param page query int false "page"
// @Param page_size query int false "size"
// @Success 200 {object} map[string]any
// @Router /interview-tracker/authen/cards/{id}/comments [get]
func (h *CardHandler) ListComments(c *gin.Context) {
	id := c.Param("id")
	page := atoiDefault(c.Query("page"), 1)
	size := atoiDefault(c.Query("page_size"), 10)
	items, total, err := h.uc.ListComments(id, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "page": page, "page_size": size})
}

// @Summary Add progress log
// @Tags progress
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "card id"
// @Param request body card_models.AddProgressReq true "progress message"
// @Success 200 {object} map[string]any
// @Router /interview-tracker/authen/cards/{id}/progress [post]
func (h *CardHandler) AddProgress(c *gin.Context) {
	id := c.Param("id")
	var req card_models.AddProgressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	session := middleware.GetSession(c)
	if err := h.uc.AddProgress(session.UserID, id, req.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// @Summary List progress logs
// @Tags progress
// @Security BearerAuth
// @Produce json
// @Param id path string true "card id"
// @Param page query int false "page"
// @Param page_size query int false "size"
// @Success 200 {object} map[string]any
// @Router /interview-tracker/authen/cards/{id}/progress [get]
func (h *CardHandler) ListProgress(c *gin.Context) {
	id := c.Param("id")
	page := atoiDefault(c.Query("page"), 1)
	size := atoiDefault(c.Query("page_size"), 10)
	items, total, err := h.uc.ListProgress(id, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "page": page, "page_size": size})
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}
