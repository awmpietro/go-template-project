package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nuhorizon/go-project-template/services/template/internal/models"
	"github.com/nuhorizon/go-project-template/services/template/internal/usecases"

	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
	ExchangeToken(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authUseCase usecases.AuthUseCase
	validator   *validator.Validate
}

func NewAuthHandler(authUC usecases.AuthUseCase) AuthHandler {
	return &authHandler{
		authUseCase: authUC,
		validator:   validator.New(),
	}
}

// Login godoc
// @Summary Realiza o login e sincronização do usuário com Firebase
// @Description Recebe o token do Firebase, valida, sincroniza e retorna o token da aplicação
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Firebase Token"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/login [post]
func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := a.validator.Struct(req); err != nil {
		httpError(w, http.StatusBadRequest, "validation error")
		return
	}

	user, token, err := a.authUseCase.LoginOrRegister(r.Context(), req.FirebaseToken)
	if err != nil {
		httpError(w, http.StatusUnauthorized, fmt.Sprintf("invalid token: %v", err.Error()))
		return
	}

	response := models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			PictureURL: user.PictureURL,
			PlanType:   user.PlanType,
		},
	}

	httpSuccess(w, http.StatusOK, response)
}

// Register - (Opcional) Se suportar registro direto sem Firebase
func (a *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	httpError(w, http.StatusNotImplemented, "not implemented")
}

// ResetPassword godoc
// @Summary Envia o e-mail de recuperação de senha via Firebase
// @Description Recebe o e-mail e dispara o fluxo de reset de senha
// @Tags Auth
// @Accept json
// @Produce json
// @Param resetPassword body models.ResetPasswordRequest true "Email para reset de senha"
// @Success 204 "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/reset-password [post]
func (a *authHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := a.validator.Struct(req); err != nil {
		httpError(w, http.StatusBadRequest, "validation error")
		return
	}

	if err := a.authUseCase.ResetPassword(r.Context(), req.Email); err != nil {
		httpError(w, http.StatusInternalServerError, fmt.Sprintf("reset password failed: %v", err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ExchangeToken godoc
// @Summary (Opcional) Troca o token de login por um novo token da aplicação
// @Description Recebe um token de refresh ou de terceiro e retorna o token da aplicação
// @Tags Auth
// @Accept json
// @Produce json
// @Param exchange body models.ExchangeTokenRequest true "Token para troca"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/exchange-token [post]
func (a *authHandler) ExchangeToken(w http.ResponseWriter, r *http.Request) {
	httpError(w, http.StatusNotImplemented, "not implemented")
}
