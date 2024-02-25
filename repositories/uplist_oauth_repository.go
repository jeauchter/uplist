package repositories

import (
	"github.com/jeremyauchter/uplist/models"
	"gorm.io/gorm"
)

type UplistOauthRepository struct {
	db *gorm.DB
}

func NewUplistOauthRepository(db *gorm.DB) *UplistOauthRepository {
	return &UplistOauthRepository{
		db: db,
	}
}

func (r *UplistOauthRepository) Create(oauth *models.UplistOauth) error {
	return r.db.Create(oauth).Error
}

func (r *UplistOauthRepository) Read(id uint) (*models.UplistOauth, error) {
	var oauth models.UplistOauth
	err := r.db.First(&oauth, id).Error
	if err != nil {
		return nil, err
	}
	return &oauth, nil
}

func (r *UplistOauthRepository) Update(oauth *models.UplistOauth) error {
	return r.db.Save(oauth).Error
}

func (r *UplistOauthRepository) Delete(id uint) error {
	return r.db.Delete(&models.UplistOauth{}, id).Error
}

func (r *UplistOauthRepository) GetAll() ([]models.UplistOauth, error) {
	var oauths []models.UplistOauth
	err := r.db.Find(&oauths).Error
	if err != nil {
		return nil, err
	}
	return oauths, nil
}

func (r *UplistOauthRepository) UpdateByID(id uint, oauth *models.UplistOauth) error {
	return r.db.Model(&models.UplistOauth{}).Where("id = ?", id).Updates(oauth).Error
}
