package dto

type UserDto struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	IsActive bool
}
