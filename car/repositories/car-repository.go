package repositories

import (
	"C_CRUD/car/repositories/models"
	"gorm.io/gorm"
)

func NewCarRepository(db *gorm.DB) *CarRepository {
	return &CarRepository{db: db}
}

type CarRepository struct {
	db *gorm.DB
}

func (r CarRepository) List() ([]models.Car, error) {
	var cars []models.Car

	if err := r.db.Find(&cars).Error; err != nil {
		return nil, err
	}

	return cars, nil
}

func (r CarRepository) Create(modelName, modelType, transmission, engine, horsePower string) (models.Car, error) {
	car := models.Car{
		ModelName:    modelName,
		Type:         modelType,
		Transmission: transmission,
		Engine:       engine,
		HorsePower:   horsePower,
	}

	if err := r.db.Create(&car).Error; err != nil {
		return models.Car{}, err
	}

	return car, nil
}

func (r CarRepository) Get(id string) (models.Car, error) {
	var car models.Car
	if err := r.db.Where("id = ?", id).First(&car).Error; err != nil {
		return models.Car{}, err
	}

	return car, nil
}

func (r CarRepository) Update(id string, modelName, modelType, transmission, engine, horsePower string) (models.Car, error) {
	car := models.Car{
		ModelName:    modelName,
		Type:         modelType,
		Transmission: transmission,
		Engine:       engine,
		HorsePower:   horsePower,
	}

	if err := r.db.Where("id = ?", id).Updates(&car).Error; err != nil {
		return models.Car{}, err
	}

	return car, nil
}

func (r CarRepository) Delete(id string) error {
	if err := r.db.Where("id = ?", id).Delete(models.Car{}).Error; err != nil {
		return err
	}

	return nil
}
