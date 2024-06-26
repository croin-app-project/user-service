package repositories

import (
	"github.com/croin-app-project/user-service/internal/domain"

	"gorm.io/gorm"
)

type UserGormRepository struct {
	db *gorm.DB
}

// NewUserGormRepository creates a new instance of userGormRepository
func NewUserGormRepository(db *gorm.DB) domain.IUserRepository {
	return &UserGormRepository{db}
}

// Create inserts a new user record into the database
func (r *UserGormRepository) Create(user *domain.User) error {
	return r.db.Model(&domain.User{}).Create(user).Error
}

// FindByID retrieves a user record by its ID from the database
func (r *UserGormRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.Model(&domain.User{}).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByCredential retrieves a user record by its username and password a from the database
func (r *UserGormRepository) FindByCredential(username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Model(&domain.User{}).Where(&domain.User{Username: username, IsActive: true}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserGormRepository) IsExistsByUsername(username string) (bool, error) {
	if err := r.db.Model(&domain.User{}).Where(&domain.User{Username: username, IsActive: true}).Select("username").First(&username).Error; err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// Update updates an existing user record in the database
func (r *UserGormRepository) Update(user *domain.User) error {
	return r.db.Model(&domain.User{}).Save(user).Error
}

// Delete removes a user record from the database by its ID
func (r *UserGormRepository) Delete(id uint) error {
	return r.db.Model(&domain.User{}).Delete(&domain.User{}, id).Error
}

// FindAll retrieves all user records from the database
func (r *UserGormRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Model(&domain.User{}).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
