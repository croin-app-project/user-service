package dto

type RegisterDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
}

type RegisterPresenter struct {
	Message string `json:"message"`
}

type LoginDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginPresenter struct {
	AccessToken string `json:"accessToken"`
}
