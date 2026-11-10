package controller

import (
	"Notifications/src/internal/domain"
	"Notifications/src/internal/service"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthenticationService
}

func NewAuthHandler(authService *service.AuthenticationService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequestDTO struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponseDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
}

func toUserResponseDTO(u *domain.User) userResponseDTO {
	return userResponseDTO{
		ID:       u.ID.String(),
		Name:     u.Name,
		LastName: u.LastName,
		Email:    u.Email,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto registerRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.authService.Register(r.Context(), service.RegisterRequest{
		Name:     dto.Name,
		LastName: dto.LastName,
		Email:    dto.Email,
		Password: dto.Password,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toUserResponseDTO(user))
}

type loginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponseDTO struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var dto loginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.authService.Authenticate(r.Context(), dto.Email, dto.Password)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, loginResponseDTO{Token: token})
}
