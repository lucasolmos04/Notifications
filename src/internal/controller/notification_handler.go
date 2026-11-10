package controller

import (
	"Notifications/src/internal/domain"
	"Notifications/src/internal/service"
	"encoding/json"
	"net/http"
	"uuid"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

type createNotificationRequestDTO struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Channel string `json:"channel"`
}

type updateNotificationRequestDTO struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Channel string `json:"channel"`
}

type notificationResponseDTO struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Channel string `json:"channel"`
	Status  string `json:"status"`
}

func toNotificationResponseDTO(n *domain.Notification) notificationResponseDTO {
	return notificationResponseDTO{
		ID:      n.ID.String(),
		Title:   n.Title,
		Content: n.Content,
		Channel: string(n.Channel),
		Status:  string(n.Status),
	}
}

func (h *NotificationHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var dto createNotificationRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	notification, err := h.notificationService.Create(r.Context(), service.CreateNotificationRequest{
		UserID:  userID,
		Title:   dto.Title,
		Content: dto.Content,
		Channel: domain.Channel(dto.Channel),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toNotificationResponseDTO(notification))
}

func (h *NotificationHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	notificationID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid notification id")
		return
	}

	var dto updateNotificationRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	notification, err := h.notificationService.Update(r.Context(), userID, notificationID, service.UpdateNotificationRequest{
		Title:   dto.Title,
		Content: dto.Content,
		Channel: domain.Channel(dto.Channel),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toNotificationResponseDTO(notification))
}

func (h *NotificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	notificationID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid notification id")
		return
	}

	if err := h.notificationService.Delete(r.Context(), userID, notificationID); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	notifications, err := h.notificationService.ListByUser(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response := make([]notificationResponseDTO, 0, len(notifications))
	for _, n := range notifications {
		response = append(response, toNotificationResponseDTO(&n))
	}

	writeJSON(w, http.StatusOK, response)
}
