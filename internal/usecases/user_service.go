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

	// Create user
	user := &domain.User{Username: u.Username, Email: u.Email, PasswordHash: u.PasswordHash, IsActive: true}
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

func (impl *UserServiceImpl) GetAllUsers() ([]domain.User, error) {
	users, err := impl._userRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (impl *UserServiceImpl) UpdateUser(u domain.User) error {
	user, err := impl._userRepository.FindByID(u.UserID)
	if err != nil {
		return err
	}

	user.Username = u.Username
	user.Email = u.Email
	user.PasswordHash = u.PasswordHash
	user.IsActive = u.IsActive

	if err := impl._userRepository.Update(user); err != nil {
		return err
	}

	return nil
}
