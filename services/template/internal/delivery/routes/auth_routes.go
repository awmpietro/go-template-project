package routes

import (
	handlers "github.com/nuhorizon/go-project-template/services/template/internal/delivery/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router, h handlers.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)                  // POST /auth/login - Recebe Firebase token e faz login/sync
		r.Post("/register", h.Register)            // POST /auth/register - (Opcional) registro direto
		r.Post("/reset-password", h.ResetPassword) // POST /auth/reset-password - Esqueci minha senha
		r.Post("/exchange-token", h.ExchangeToken) // POST /auth/exchange-token - (opcional) fluxo para trocar token
	})
}
