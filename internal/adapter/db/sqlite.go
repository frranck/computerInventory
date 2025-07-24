package db

import (
	"computerInventory/internal/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteRepo struct {
	db *gorm.DB
}

func NewSQLiteRepo(path string) (*SQLiteRepo, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&domain.Computer{})
	return &SQLiteRepo{db: db}, err
}

func (r *SQLiteRepo) Create(c *domain.Computer) error {
	return r.db.Create(c).Error
}

func (r *SQLiteRepo) Get(mac string) (*domain.Computer, error) {
	var comp domain.Computer
	err := r.db.First(&comp, "mac_address = ?", mac).Error
	return &comp, err
}

func (r *SQLiteRepo) Update(c *domain.Computer) error {
	return r.db.Save(c).Error
}

func (r *SQLiteRepo) Delete(mac string) error {
	return r.db.Delete(&domain.Computer{}, "mac_address = ?", mac).Error
}

func (r *SQLiteRepo) GetAll() ([]domain.Computer, error) {
	var comps []domain.Computer
	err := r.db.Find(&comps).Error
	return comps, err
}

func (r *SQLiteRepo) GetByEmployee(abbr string) ([]domain.Computer, error) {
	var comps []domain.Computer
	err := r.db.Where("employee_abbreviation = ?", abbr).Find(&comps).Error
	return comps, err
}
