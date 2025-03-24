package models

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	PictureURL string `json:"picture_url,omitempty"`
	PlanType   string `json:"plan_type"`
}
