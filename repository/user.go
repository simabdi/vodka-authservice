package repository

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/simabdi/vodka-authservice/helper"
	"github.com/simabdi/vodka-authservice/models"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserRepository interface {
	Store(ctx *fiber.Ctx, input models.User) (models.User, error)
	Update(input models.User) (models.User, error)
	UpdateColumn(input models.User, column string, value string) error
	GetByEmail(email string) (models.User, error)
	GetById(userId int) (models.User, error)
	GetByUuid(uuid string) (models.User, error)
	GetAll(ctx *fiber.Ctx) ([]models.User, error)
	GetByRef(refID uint, refType string) (models.User, error)
	Transaction(ctx *fiber.Ctx, fn func(repo UserRepository) error) error
}

type repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) withTx(tx *gorm.DB) UserRepository {
	return &repository{
		db: tx,
	}
}

func statusActive(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", "active")
}

func (r *repository) Transaction(ctx *fiber.Ctx, fn func(repo UserRepository) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	repo := r.withTx(tx)
	err := fn(repo)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *repository) Store(ctx *fiber.Ctx, input models.User) (models.User, error) {
	err := r.db.Create(&input).Error
	if err != nil {
		return input, err
	}

	return input, nil
}

func (r *repository) Update(input models.User) (models.User, error) {
	err := r.db.Save(input).Error
	if err != nil {
		return input, err
	}

	return input, nil
}

func (r *repository) UpdateColumn(input models.User, column string, value string) error {
	err := r.db.Model(&input).UpdateColumn(column, value).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByEmail(email string) (models.User, error) {
	var user models.User

	err := r.db.Where("email = ?", email).First(&user).Error
	r.db.Logger = logger.Default.LogMode(logger.Info)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) GetById(userId int) (models.User, error) {
	var user models.User

	err := r.db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) GetByUuid(uuid string) (models.User, error) {
	var user models.User

	err := r.db.Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) GetAll(ctx *fiber.Ctx) ([]models.User, error) {
	var user []models.User
	var count int64

	err := r.db.Find(&user, "role != ?", "ADMIN").Error
	if err != nil {
		return user, err
	}
	r.db.Scopes(statusActive, helper.Paginate(ctx)).Model(&user).Count(&count)
	helper.TotalRecord = count

	r.db.Logger = logger.Default.LogMode(logger.Info)
	return user, nil
}

func (r *repository) GetByRef(refID uint, refType string) (models.User, error) {
	var user models.User
	err := r.db.Where("ref_id = ? AND ref_type = ?", refID, refType).Find(&user).Error
	if err != nil {
		return user, err
	}
	errors.Is(err, gorm.ErrRecordNotFound)

	r.db = r.db.Debug()
	return user, nil
}
