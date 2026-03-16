package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Novodremov/subscribe-be/internal/dto"
	"github.com/Novodremov/subscribe-be/internal/repo"
	"github.com/Novodremov/subscribe-be/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type SubscriptionHandler struct {
	svc service.ISubscriptionService
	log zerolog.Logger
}

// NewSubscriptionHandler создает новый экземпляр SubscriptionHandler с переданным сервисом подписок и логгером.
// Используется для инициализации хендлеров HTTP маршрутов.
func NewSubscriptionHandler(svc service.ISubscriptionService, logger zerolog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		svc: svc,
		log: logger,
	}
}

// CreateSubscription создаёт новую подписку.
//
// @Summary Create subscription
// @Accept json
// @Produce json
// @Param body body dto.CreateSubscriptionRequest true "Create Subscription"
// @Success 201 {object} dto.SubscriptionResponse
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /subscription [post]
func (h *SubscriptionHandler) CreateSubscription(c *fiber.Ctx) error {
	const op = "CreateSubscription"

	var req dto.CreateSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error().Err(err).Str("op", op).Msg("invalid request body")
		return NewHTTPError(err, http.StatusBadRequest, "invalid request body")
	}

	if req.ServiceName == "" {
		h.log.Error().Str("op", op).Msg("service_name is missing")
		return NewHTTPError(ErrValidation, http.StatusBadRequest, "service_name is required")
	}
	if req.Price <= 0 {
		h.log.Error().Int("price", req.Price).Str("op", op).Msg("invalid price")
		return NewHTTPError(ErrValidation, http.StatusBadRequest, "price must be > 0")
	}
	if req.UserID == uuid.Nil {
		h.log.Error().Str("op", op).Msg("user_id is missing")
		return NewHTTPError(ErrValidation, http.StatusBadRequest, "user_id is required")
	}

	domainSub, err := MapCreateDTOToCreateDomain(req)
	if err != nil {
		h.log.Error().Err(err).Str("op", op).Msg("failed to map DTO to domain")
		return NewHTTPError(err, http.StatusBadRequest, "invalid subscription data")
	}

	if domainSub.EndDate != nil && domainSub.EndDate.Before(domainSub.StartDate) {
		h.log.Error().
			Time("start_date", domainSub.StartDate).
			Time("end_date", *domainSub.EndDate).
			Str("op", op).
			Msg("end_date cannot be before start_date")
		return NewHTTPError(ErrValidation, http.StatusBadRequest, "end_date cannot be before start_date")
	}

	created, err := h.svc.CreateSubscription(c.Context(), &domainSub)
	if err != nil {
		h.log.Error().Err(err).Str("op", op).Msg("failed to create subscription")
		return NewHTTPError(err, http.StatusInternalServerError, "failed to create subscription")
	}

	h.log.Info().
		Str("op", op).
		Str("subscription_id", created.ID.String()).
		Str("user_id", created.UserID.String()).
		Msg("subscription successfully created")

	return c.Status(http.StatusCreated).JSON(MapDomainToResponse(*created))
}

// GetSubscription извлекает подписку по ID.
//
// @Summary Get subscription
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /subscription/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *fiber.Ctx) error {
	const op = "GetSubscription"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("op", op).Msg("invalid subscription ID")
		return NewHTTPError(err, http.StatusBadRequest, "invalid subscription ID")
	}

	sub, err := h.svc.GetSubscription(c.Context(), id)

	if err != nil {
		if errors.Is(err, repo.ErrSubscriptionNotFound) {
			h.log.Error().Str("op", op).Str("subscription_id", id.String()).Msg("subscription not found")
			return NewHTTPError(err, http.StatusNotFound, "subscription not found")
		}
		h.log.Error().Err(err).Str("op", op).Msg("failed to get subscription")
		return NewHTTPError(err, http.StatusInternalServerError, "failed to get subscription")
	}

	h.log.Info().
		Str("op", op).
		Str("subscription_id", sub.ID.String()).
		Msg("subscription retrieved successfully")

	return c.JSON(MapDomainToResponse(*sub))
}

// UpdateSubscription обновляет данные по подписке.
//
// @Summary Update subscription
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param body body dto.UpdateSubscriptionRequest true "Update Subscription"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /subscription/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *fiber.Ctx) error {
	const op = "UpdateSubscription"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("op", op).Msg("invalid subscription ID")
		return NewHTTPError(err, http.StatusBadRequest, "invalid subscription ID")
	}

	var req dto.UpdateSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error().Err(err).Str("op", op).Msg("failed to parse request body")
		return NewHTTPError(err, http.StatusBadRequest, "invalid request body")
	}

	if req.ServiceName == nil && req.Price == nil && req.StartDate == nil && req.EndDate == nil {
		h.log.Error().Str("op", op).Msg("no fields provided for update")
		return NewHTTPError(ErrNoFieldsToUpdate, http.StatusBadRequest, "at least one field must be provided for update")
	}

	domainSub, err := MapUpdateDTOToUpdateDomain(req)
	if err != nil {
		h.log.Error().Err(err).Str("op", op).Msg("failed to map DTO to domain")
		return NewHTTPError(err, http.StatusBadRequest, "invalid subscription data")
	}

	domainSub.ID = id

	if domainSub.StartDate != nil && domainSub.EndDate != nil && domainSub.EndDate.Before(*domainSub.StartDate) {
		h.log.Error().
			Time("start_date", *domainSub.StartDate).
			Time("end_date", *domainSub.EndDate).
			Str("op", op).
			Msg("end_date cannot be before start_date")
		return NewHTTPError(ErrValidation, http.StatusBadRequest, "end_date cannot be before start_date")
	}
	
	updated, err := h.svc.UpdateSubscription(c.Context(), &domainSub)
	if err != nil {
		if errors.Is(err, repo.ErrSubscriptionNotFound) {
			h.log.Error().Str("op", op).Str("subscription_id", id.String()).Msg("subscription not found")
			return NewHTTPError(err, http.StatusNotFound, "subscription not found")
		}
		h.log.Error().Err(err).Str("op", op).Msg("failed to update subscription")
		return NewHTTPError(err, http.StatusInternalServerError, "failed to update subscription")
	}
	h.log.Info().
		Str("op", op).
		Str("subscription_id", updated.ID.String()).
		Msg("subscription updated successfully")

	return c.JSON(MapDomainToResponse(*updated))
}

// DeleteSubscription удаляет подписку по ID.
//
// @Summary Delete subscription
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /subscription/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *fiber.Ctx) error {
	const op = "DeleteSubscription"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("op", op).Msg("invalid subscription ID")
		return NewHTTPError(err, http.StatusBadRequest, "invalid subscription ID")
	}

	if err := h.svc.DeleteSubscription(c.Context(), id); err != nil {
		if errors.Is(err, repo.ErrSubscriptionNotFound) {
			h.log.Error().Str("op", op).Str("subscription_id", id.String()).Msg("subscription not found")
			return NewHTTPError(err, http.StatusNotFound, "subscription not found")
		}
		h.log.Error().Err(err).Str("op", op).Msg("failed to delete subscription")
		return NewHTTPError(err, http.StatusInternalServerError, "failed to delete subscription")
	}
	h.log.Info().
		Str("op", op).
		Str("subscription_id", id.String()).
		Msg("subscription deleted successfully")

	return c.SendStatus(http.StatusNoContent)
}

// ListSubscriptions возвращает список подписок с учётом пагинации (лимит, офсет) и их общее количество.
//
// @Summary List of subscriptions
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} dto.ListSubscriptionsResponse
// @Failure 500 {object} HTTPError
// @Router /subscription [get]
func (h *SubscriptionHandler) ListSubscriptions(c *fiber.Ctx) error {
	const op = "ListSubscriptions"

	limit, offset, err := parseLimitOffset(c)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("op", op).
			Msg("invalid pagination parameters")
		return NewHTTPError(err, http.StatusBadRequest, "invalid pagination parameters")
	}

	items, total, err := h.svc.ListSubscriptions(c.Context(), limit, offset)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("op", op).
			Msg("failed to list subscriptions")
		return NewHTTPError(err, http.StatusInternalServerError, "failed to list subscriptions")
	}

	h.log.Info().
		Int("count", len(items)).
		Int("total_count", total).
		Str("op", op).
		Msg("subscriptions listed successfully")

	return c.JSON(MapDomainSubscriptionsToDTO(items, total))
}

// SubscriptionsTotalCost вычисляет суммарную стоимость подписок для указанного пользователя и/или сервиса за заданный период.
// Параметры startDate и endDate могут быть nil, тогда период не учитывается.
//
// @Summary Get total cost of subscriptions with filters
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service Name"
// @Param start_date query string false "Start date (DD-MM-YYYY)"
// @Param end_date query string false "End date (DD-MM-YYYY)"
// @Success 200 {object} dto.SubscriptionsTotalCostResponse
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /subscription/total-cost [get]
func (h *SubscriptionHandler) SubscriptionsTotalCost(c *fiber.Ctx) error {
	const op = "SubscriptionsTotalCost"

	var userID *uuid.UUID
	if s := c.Query("user_id"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			h.log.Error().
				Err(err).
				Str("op", op).
				Str("user_id", s).
				Msg("invalid user ID")
			return NewHTTPError(err, http.StatusBadRequest, "invalid user ID")
		}
		userID = &id
	}

	var serviceName *string
	if s := c.Query("service_name"); s != "" {
		serviceName = &s
	}

	var startDate *time.Time
	if s := c.Query("start_date"); s != "" {
		parsed, err := time.Parse(DateLayout, s)
		if err != nil {
			return NewHTTPError(err, http.StatusBadRequest, "invalid start_date format, expected DD-MM-YYYY")
		}
		startDate = &parsed
	}

	var endDate *time.Time
	if s := c.Query("end_date"); s != "" {
		parsed, err := time.Parse(DateLayout, s)
		if err != nil {
			return NewHTTPError(err, http.StatusBadRequest, "invalid end_date format, expected DD-MM-YYYY")
		}
		endDate = &parsed
	}

	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return NewHTTPError(
			fmt.Errorf("start_date after end_date"),
			http.StatusBadRequest,
			"start_date cannot be later than end_date",
		)
	}

	total, err := h.svc.SubscriptionsTotalCost(c.Context(), userID, serviceName, startDate, endDate)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("op", op).
			Msg("failed to list subscriptions with filters")
		return NewHTTPError(err, http.StatusInternalServerError, "failed to list subscriptions with filters")
	}

	h.log.Info().
		Int64("total_cost", total).
		Str("op", op).
		Msg("total cost of filtered subscriptions calculated successfully")

	return c.JSON(dto.SubscriptionsTotalCostResponse{
		TotalCost: total,
	})
}

func parseLimitOffset(c *fiber.Ctx) (int, int, error) {
	limit := DefaultLimit
	offset := DefaultOffset

	if l := c.Query("limit"); l != "" {
		if _, err := fmt.Sscan(l, &limit); err != nil {
			return 0, 0, fmt.Errorf("invalid limit: %w", err)
		}
		if limit < MinLimit || limit > MaxLimit {
			return 0, 0, fmt.Errorf("limit must be between %d and %d", MinLimit, MaxLimit)
		}
	}

	if o := c.Query("offset"); o != "" {
		if _, err := fmt.Sscan(o, &offset); err != nil {
			return 0, 0, fmt.Errorf("invalid offset: %w", err)
		}
		if offset < 0 {
			return 0, 0, fmt.Errorf("offset cannot be negative")
		}
	}

	return limit, offset, nil
}
