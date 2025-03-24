package domain

import "time"

type User struct {
	ID           string
	FirebaseUID  string
	Email        string
	Name         string
	PictureURL   string
	PlanType     string
	PremiumSince *time.Time
	PlanExpiry   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
