package postgres

import (
	"errors"
	"todo-app/domain"
	"todo-app/pkg/clients"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Save(user *domain.UserCreate) error {
	if err := r.db.Create(&user).Error; err != nil {
		return clients.ErrDB(err)
	}
	return nil
}

func (r *userRepo) GetUser(conditions map[string]any) (*domain.User, error) {
	var user domain.User

	if err := r.db.Where(conditions).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, clients.ErrRecordNotFound
		}
		return nil, clients.ErrDB(err)
	}

	return &user, nil
}


func (r *userRepo) GetAll() ([]*domain.User, error) {
	var users []*domain.User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, clients.ErrDB(err)
	}

	return users, nil
}


func (r *userRepo) UpdateUser(user *domain.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return clients.ErrDB(err)
	}

	return nil
}


func (r *userRepo) DeleteUser(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&domain.User{}).Error; err != nil {
		return clients.ErrDB(err)
	}

	return nil
}
