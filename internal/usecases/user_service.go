package usecases

import (
	helpers "github.com/croin-app-project/package/pkg/utils/helpers"
	"github.com/croin-app-project/user-service/internal/domain"
	"github.com/croin-app-project/user-service/internal/usecases/iservices"
)

type UserServiceImpl struct {
	_userRepository domain.IUserRepository
}

func NewUserService(userRepository domain.IUserRepository) iservices.IUserService {
	return &UserServiceImpl{
		_userRepository: userRepository,
	}
}

func (impl *UserServiceImpl) IsAlreadyExistsByUsername(useranme string) (bool, error) {
	isExists, err := impl._userRepository.IsExistsByUsername(useranme)
	if err != nil {
		return false, err
	}
	return isExists, nil
}

func (impl *UserServiceImpl) SaveNewUser(u domain.User) error {
	// Hash the password
	hashedPassword, err := helpers.HashPassword(u.PasswordHash)
	if err != nil {
		return err
	}

	// Create user
	user := &domain.User{Username: u.Username, Email: u.Email, PasswordHash: hashedPassword, IsActive: true}
	if err := impl._userRepository.Create(user); err != nil {
		return err
	}

	return nil
}

func (impl *UserServiceImpl) VerifyUsernamePassword(username string, password string) (*domain.User, error) {
	user, err := impl._userRepository.FindByCredential(username)
	if err != nil {
		return nil, err
	}

	if helpers.ComparePassword(user.PasswordHash, password) {
		return user, nil
	} else {
		return nil, nil
	}
}
