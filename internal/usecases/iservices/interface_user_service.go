package iservices

import "github.com/croin-app-project/user-service/internal/domain"

type IUserService interface {
	IsAlreadyExistsByUsername(useranme string) (bool, error)
	SaveNewUser(user domain.User) error
	VerifyUsernamePassword(username string, password string) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	UpdateUser(user domain.User) error
}
