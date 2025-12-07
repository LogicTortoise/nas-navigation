package services

import (
	"nas-navigation/database"
	"nas-navigation/models"

	"gorm.io/gorm"
)

type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := database.DB.Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

func (s *CategoryService) GetAllWithLinks() ([]models.Category, error) {
	var categories []models.Category
	err := database.DB.Preload("Links", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

func (s *CategoryService) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := database.DB.Preload("Links").First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) Create(req *models.CategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:      req.Name,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
	}
	err := database.DB.Create(category).Error
	return category, err
}

func (s *CategoryService) Update(id uint, req *models.CategoryRequest) (*models.Category, error) {
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Icon = req.Icon
	category.SortOrder = req.SortOrder

	err := database.DB.Save(&category).Error
	return &category, err
}

func (s *CategoryService) Delete(id uint) error {
	// 先删除该分类下的所有链接
	if err := database.DB.Where("category_id = ?", id).Delete(&models.Link{}).Error; err != nil {
		return err
	}
	return database.DB.Delete(&models.Category{}, id).Error
}
