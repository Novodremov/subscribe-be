package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	subhttp "github.com/Novodremov/subscribe-be/internal/delivery/http"
	"github.com/Novodremov/subscribe-be/internal/delivery/http/handler"
	"github.com/Novodremov/subscribe-be/internal/domain"
	"github.com/Novodremov/subscribe-be/internal/dto"
	"github.com/Novodremov/subscribe-be/internal/repo"
	mock_service "github.com/Novodremov/subscribe-be/internal/service/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionHandler_CreateSubscription_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Post("/subscriptions", h.CreateSubscription)

	now := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	reqDTO := dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      uuid.New(),
		StartDate:   "01-2026",
		EndDate:     "12-2026",
	}

	expectedDomain, err := handler.MapCreateDTOToCreateDomain(reqDTO)
	require.NoError(t, err)

	created := &domain.Subscription{
		ID:          uuid.New(),
		ServiceName: expectedDomain.ServiceName,
		Price:       expectedDomain.Price,
		UserID:      expectedDomain.UserID,
		StartDate:   expectedDomain.StartDate,
		EndDate:     expectedDomain.EndDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockSvc.EXPECT().
		CreateSubscription(gomock.Any(), gomock.Eq(&expectedDomain)).
		Return(created, nil).
		Times(1)

	body, err := json.Marshal(reqDTO)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var actual dto.SubscriptionResponse
	err = json.NewDecoder(resp.Body).Decode(&actual)
	require.NoError(t, err)

	expected := handler.MapDomainToResponse(*created)
	assert.Equal(t, expected, actual)
}

func TestSubscriptionHandler_CreateSubscription_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Post("/subscriptions", h.CreateSubscription)

	invalidBody := []byte(`{"service_name": "Netflix", "price": 599,`) // некорректный JSON

	req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(invalidBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_CreateSubscription_EmptyServiceName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Post("/subscriptions", h.CreateSubscription)

	reqDTO := dto.CreateSubscriptionRequest{
		ServiceName: "",
		Price:       599,
		UserID:      uuid.New(),
		StartDate:   "01-2026",
		EndDate:     "12-2026",
	}

	body, err := json.Marshal(reqDTO)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_CreateSubscription_InvalidPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Post("/subscriptions", h.CreateSubscription)

	reqDTO := dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       0, // <= 0
		UserID:      uuid.New(),
		StartDate:   "01-2026",
		EndDate:     "12-2026",
	}

	body, err := json.Marshal(reqDTO)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_CreateSubscription_EmptyUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Post("/subscriptions", h.CreateSubscription)

	reqDTO := dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      uuid.Nil,
		StartDate:   "01-2026",
		EndDate:     "12-2026",
	}

	body, err := json.Marshal(reqDTO)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_CreateSubscription_InvalidUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Post("/subscriptions", h.CreateSubscription)

	reqDTO := dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      uuid.Nil, // невалидный user_id
		StartDate:   "01-2026",
		EndDate:     "12-2026",
	}

	body, err := json.Marshal(reqDTO)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_CreateSubscription_InvalidStartDateFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Post("/subscriptions", h.CreateSubscription)

	reqDTO := dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      uuid.New(),
		StartDate:   "2026-01", // неправильный формат
		EndDate:     "12-2026",
	}

	body, err := json.Marshal(reqDTO)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_CreateSubscription_InvalidEndDateFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Post("/subscriptions", h.CreateSubscription)

	reqDTO := dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      uuid.New(),
		StartDate:   "01-2026",
		EndDate:     "2026-12", // неправильный формат
	}

	body, err := json.Marshal(reqDTO)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_CreateSubscription_EndDateBeforeStartDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Post("/subscriptions", h.CreateSubscription)

	reqDTO := dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      uuid.New(),
		StartDate:   "12-2026",
		EndDate:     "01-2026", // конец раньше начала
	}

	body, err := json.Marshal(reqDTO)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_CreateSubscription_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Post("/subscriptions", h.CreateSubscription)

	reqDTO := dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      uuid.New(),
		StartDate:   "01-2026",
		EndDate:     "12-2026",
	}

	mockSvc.EXPECT().
		CreateSubscription(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("db exploded")).
		Times(1)

	body, err := json.Marshal(reqDTO)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestSubscriptionHandler_GetSubscription_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Get("/subscriptions/:id", h.GetSubscription)

	id := uuid.New()

	sub := &domain.Subscription{
		ID:          id,
		ServiceName: "Netflix",
		Price:       599,
		UserID:      uuid.New(),
		StartDate:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     nil,
	}

	mockSvc.EXPECT().
		GetSubscription(gomock.Any(), id).
		Return(sub, nil).
		Times(1)

	req := httptest.NewRequest("GET", "/subscriptions/"+id.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var actual dto.SubscriptionResponse
	err = json.NewDecoder(resp.Body).Decode(&actual)
	require.NoError(t, err)

	expected := handler.MapDomainToResponse(*sub)

	assert.Equal(t, expected, actual)
}

func TestSubscriptionHandler_GetSubscription_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Get("/subscriptions/:id", h.GetSubscription)

	req := httptest.NewRequest("GET", "/subscriptions/not-a-uuid", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_GetSubscription_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Get("/subscriptions/:id", h.GetSubscription)

	id := uuid.New()

	mockSvc.EXPECT().
		GetSubscription(gomock.Any(), id).
		Return(nil, repo.ErrSubscriptionNotFound).
		Times(1)

	req := httptest.NewRequest("GET", "/subscriptions/"+id.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestSubscriptionHandler_GetSubscription_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Get("/subscriptions/:id", h.GetSubscription)

	id := uuid.New()

	mockSvc.EXPECT().
		GetSubscription(gomock.Any(), id).
		Return(nil, errors.New("db error")).
		Times(1)

	req := httptest.NewRequest("GET", "/subscriptions/"+id.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestSubscriptionHandler_UpdateSubscription_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Put("/subscriptions/:id", h.UpdateSubscription)

	now := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	subID := uuid.New()
	reqDTO := dto.UpdateSubscriptionRequest{
		ServiceName: ptrString("Netflix Premium"),
		Price:       ptrInt(699),
	}

	expectedDomain, err := handler.MapUpdateDTOToUpdateDomain(reqDTO)
	require.NoError(t, err)
	expectedDomain.ID = subID

	updated := &domain.Subscription{
		ID:          expectedDomain.ID,
		ServiceName: derefString(expectedDomain.ServiceName),
		Price:       derefInt(expectedDomain.Price),
		UserID:      uuid.New(), // просто для примера, реально можно фиксировать
		StartDate:   now,
		EndDate:     nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockSvc.EXPECT().
		UpdateSubscription(gomock.Any(), gomock.Eq(&expectedDomain)).
		Return(updated, nil).
		Times(1)

	body, err := json.Marshal(reqDTO)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "/subscriptions/"+subID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var actual dto.SubscriptionResponse
	err = json.NewDecoder(resp.Body).Decode(&actual)
	require.NoError(t, err)

	expected := handler.MapDomainToResponse(*updated)
	assert.Equal(t, expected, actual)
}

func TestSubscriptionHandler_UpdateSubscription_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Put("/subscriptions/:id", h.UpdateSubscription)

	req := httptest.NewRequest("PUT", "/subscriptions/invalid-uuid", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_UpdateSubscription_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Put("/subscriptions/:id", h.UpdateSubscription)

	req := httptest.NewRequest("PUT", "/subscriptions/"+uuid.New().String(), bytes.NewReader([]byte(`invalid-json`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_UpdateSubscription_NoFieldsProvided(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Put("/subscriptions/:id", h.UpdateSubscription)

	reqDTO := dto.UpdateSubscriptionRequest{} // все поля nil
	body, _ := json.Marshal(reqDTO)

	req := httptest.NewRequest("PUT", "/subscriptions/"+uuid.New().String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_UpdateSubscription_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Put("/subscriptions/:id", h.UpdateSubscription)

	subID := uuid.New()
	reqDTO := dto.UpdateSubscriptionRequest{
		ServiceName: ptrString("Netflix Premium"),
	}
	expectedDomain, _ := handler.MapUpdateDTOToUpdateDomain(reqDTO)
	expectedDomain.ID = subID

	mockSvc.EXPECT().
		UpdateSubscription(gomock.Any(), gomock.Eq(&expectedDomain)).
		Return(nil, repo.ErrSubscriptionNotFound).
		Times(1)

	body, _ := json.Marshal(reqDTO)
	req := httptest.NewRequest("PUT", "/subscriptions/"+subID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestSubscriptionHandler_UpdateSubscription_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Put("/subscriptions/:id", h.UpdateSubscription)

	subID := uuid.New()
	reqDTO := dto.UpdateSubscriptionRequest{
		ServiceName: ptrString("Netflix Premium"),
	}
	expectedDomain, _ := handler.MapUpdateDTOToUpdateDomain(reqDTO)
	expectedDomain.ID = subID

	mockSvc.EXPECT().
		UpdateSubscription(gomock.Any(), gomock.Eq(&expectedDomain)).
		Return(nil, errors.New("some db error")).
		Times(1)

	body, _ := json.Marshal(reqDTO)
	req := httptest.NewRequest("PUT", "/subscriptions/"+subID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestSubscriptionHandler_UpdateSubscription_EndDateBeforeStartDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Put("/subscriptions/:id", h.UpdateSubscription)

	subID := uuid.New()
	reqDTO := dto.UpdateSubscriptionRequest{
		StartDate: ptrString("12-2026"),
		EndDate:   ptrString("01-2026"),
	}

	body, _ := json.Marshal(reqDTO)
	req := httptest.NewRequest("PUT", "/subscriptions/"+subID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_DeleteSubscription_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Delete("/subscriptions/:id", h.DeleteSubscription)

	subID := uuid.New()

	mockSvc.EXPECT().
		DeleteSubscription(gomock.Any(), subID).
		Return(nil).
		Times(1)

	req := httptest.NewRequest("DELETE", "/subscriptions/"+subID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func TestSubscriptionHandler_DeleteSubscription_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Delete("/subscriptions/:id", h.DeleteSubscription)

	req := httptest.NewRequest("DELETE", "/subscriptions/invalid-uuid", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_DeleteSubscription_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Delete("/subscriptions/:id", h.DeleteSubscription)

	subID := uuid.New()

	mockSvc.EXPECT().
		DeleteSubscription(gomock.Any(), subID).
		Return(repo.ErrSubscriptionNotFound).
		Times(1)

	req := httptest.NewRequest("DELETE", "/subscriptions/"+subID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestSubscriptionHandler_DeleteSubscription_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Delete("/subscriptions/:id", h.DeleteSubscription)

	subID := uuid.New()

	mockSvc.EXPECT().
		DeleteSubscription(gomock.Any(), subID).
		Return(errors.New("db error")).
		Times(1)

	req := httptest.NewRequest("DELETE", "/subscriptions/"+subID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestSubscriptionHandler_ListSubscriptions_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Get("/subscriptions", h.ListSubscriptions)

	limit := 10
	offset := 0

	subs := []domain.Subscription{
		{
			ID:          uuid.New(),
			ServiceName: "Netflix",
			Price:       599,
			UserID:      uuid.New(),
			StartDate:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockSvc.EXPECT().
		ListSubscriptions(gomock.Any(), limit, offset).
		Return(subs, len(subs), nil).
		Times(1)

	req := httptest.NewRequest("GET", fmt.Sprintf("/subscriptions?limit=%d&offset=%d", limit, offset), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var actual dto.ListSubscriptionsResponse
	err = json.NewDecoder(resp.Body).Decode(&actual)
	require.NoError(t, err)

	expected := handler.MapDomainSubscriptionsToDTO(subs, len(subs))
	assert.Equal(t, expected, actual)
}

func TestSubscriptionHandler_ListSubscriptions_InvalidPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Get("/subscriptions", h.ListSubscriptions)

	req := httptest.NewRequest("GET", "/subscriptions?limit=abc&offset=-1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_ListSubscriptions_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Get("/subscriptions", h.ListSubscriptions)

	limit := 10
	offset := 0
	mockSvc.EXPECT().
		ListSubscriptions(gomock.Any(), limit, offset).
		Return(nil, 0, errors.New("database failure")).
		Times(1)

	req := httptest.NewRequest("GET", fmt.Sprintf("/subscriptions?limit=%d&offset=%d", limit, offset), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestSubscriptionHandler_ListSubscriptions_EmptyResult_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Get("/subscriptions", h.ListSubscriptions)

	limit := 10
	offset := 0
	mockSvc.EXPECT().
		ListSubscriptions(gomock.Any(), limit, offset).
		Return([]domain.Subscription{}, 0, nil).
		Times(1)

	req := httptest.NewRequest("GET", fmt.Sprintf("/subscriptions?limit=%d&offset=%d", limit, offset), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var actual dto.ListSubscriptionsResponse
	err = json.NewDecoder(resp.Body).Decode(&actual)
	require.NoError(t, err)
	assert.Empty(t, actual.Subscriptions)
	assert.Equal(t, 0, actual.TotalCount)
}

func TestSubscriptionHandler_SubscriptionsTotalCost_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	logger := zerolog.Nop()
	h := handler.NewSubscriptionHandler(mockSvc, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: subhttp.TestErrorHandler,
	})
	app.Get("/subscriptions/total-cost", h.SubscriptionsTotalCost)

	userID := uuid.New()
	serviceName := "Netflix"
	startDateStr := "01-01-2026"
	endDateStr := "31-12-2026"

	startDate, err := time.Parse(handler.DateLayout, startDateStr)
	require.NoError(t, err)
	endDate, err := time.Parse(handler.DateLayout, endDateStr)
	require.NoError(t, err)

	totalCost := int64(12345)

	mockSvc.EXPECT().
		SubscriptionsTotalCost(
			gomock.Any(),
			gomock.Eq(&userID),
			gomock.Eq(&serviceName),
			gomock.Eq(&startDate),
			gomock.Eq(&endDate),
		).
		Return(totalCost, nil).
		Times(1)

	req := httptest.NewRequest(
		"GET",
		fmt.Sprintf(
			"/subscriptions/total-cost?user_id=%s&service_name=%s&start_date=%s&end_date=%s",
			userID.String(),
			serviceName,
			startDateStr,
			endDateStr,
		),
		nil,
	)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var actual dto.SubscriptionsTotalCostResponse
	err = json.NewDecoder(resp.Body).Decode(&actual)
	require.NoError(t, err)
	assert.Equal(t, totalCost, actual.TotalCost)
}

func TestSubscriptionHandler_SubscriptionsTotalCost_InvalidUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Get("/subscriptions/total-cost", h.SubscriptionsTotalCost)

	req := httptest.NewRequest("GET",
		"/subscriptions/total-cost?user_id=invalid-uuid", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_SubscriptionsTotalCost_InvalidStartDateFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Get("/subscriptions/total-cost", h.SubscriptionsTotalCost)

	req := httptest.NewRequest("GET",
		"/subscriptions/total-cost?start_date=2026-01-01", nil) // должен быть 02-01-2006
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_SubscriptionsTotalCost_StartDateAfterEndDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Get("/subscriptions/total-cost", h.SubscriptionsTotalCost)

	req := httptest.NewRequest("GET",
		"/subscriptions/total-cost?start_date=31-12-2026&end_date=01-01-2026", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSubscriptionHandler_SubscriptionsTotalCost_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Get("/subscriptions/total-cost", h.SubscriptionsTotalCost)

	mockSvc.EXPECT().
		SubscriptionsTotalCost(gomock.Any(), gomock.Nil(), gomock.Nil(), gomock.Nil(), gomock.Nil()).
		Return(int64(0), errors.New("db failure")).
		Times(1)

	req := httptest.NewRequest("GET", "/subscriptions/total-cost", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestSubscriptionHandler_SubscriptionsTotalCost_ServiceError_WithFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockISubscriptionService(ctrl)
	h := handler.NewSubscriptionHandler(mockSvc, zerolog.Nop())

	app := fiber.New(fiber.Config{ErrorHandler: subhttp.TestErrorHandler})
	app.Get("/subscriptions/total-cost", h.SubscriptionsTotalCost)

	userID := uuid.New()
	serviceName := "Netflix"
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	mockSvc.EXPECT().
		SubscriptionsTotalCost(gomock.Any(), &userID, &serviceName, &start, &end).
		Return(int64(0), errors.New("db failure")).
		Times(1)

	req := httptest.NewRequest("GET",
		fmt.Sprintf("/subscriptions/total-cost?user_id=%s&service_name=%s&start_date=%s&end_date=%s",
			userID.String(), serviceName, start.Format(handler.DateLayout), end.Format(handler.DateLayout)),
		nil,
	)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

// вспомогательные функции для указателей
func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }
func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
func derefInt(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}
