package services

import (
	"nas-navigation/database"
	"nas-navigation/models"
)

type LinkService struct{}

func NewLinkService() *LinkService {
	return &LinkService{}
}

func (s *LinkService) GetAll() ([]models.Link, error) {
	var links []models.Link
	err := database.DB.Preload("Category").Order("sort_order ASC").Find(&links).Error
	return links, err
}

func (s *LinkService) GetByCategory(categoryID uint) ([]models.Link, error) {
	var links []models.Link
	err := database.DB.Where("category_id = ?", categoryID).
		Order("sort_order ASC").
		Find(&links).Error
	return links, err
}

func (s *LinkService) GetByID(id uint) (*models.Link, error) {
	var link models.Link
	err := database.DB.Preload("Category").First(&link, id).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (s *LinkService) Search(query string) ([]models.Link, error) {
	var links []models.Link
	searchQuery := "%" + query + "%"
	err := database.DB.Where("name LIKE ? OR description LIKE ? OR url LIKE ?",
		searchQuery, searchQuery, searchQuery).
		Preload("Category").
		Order("sort_order ASC").
		Find(&links).Error
	return links, err
}

func (s *LinkService) Create(req *models.LinkRequest) (*models.Link, error) {
	link := &models.Link{
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		URL:           req.URL,
		Icon:          req.Icon,
		Description:   req.Description,
		SortOrder:     req.SortOrder,
		RestartScript: req.RestartScript,
	}
	err := database.DB.Create(link).Error
	if err != nil {
		return nil, err
	}
	// 重新加载关联
	return s.GetByID(link.ID)
}

func (s *LinkService) Update(id uint, req *models.LinkRequest) (*models.Link, error) {
	var link models.Link
	if err := database.DB.First(&link, id).Error; err != nil {
		return nil, err
	}

	link.CategoryID = req.CategoryID
	link.Name = req.Name
	link.URL = req.URL
	link.Icon = req.Icon
	link.Description = req.Description
	link.SortOrder = req.SortOrder
	link.RestartScript = req.RestartScript

	err := database.DB.Save(&link).Error
	if err != nil {
		return nil, err
	}
	return s.GetByID(link.ID)
}

func (s *LinkService) Delete(id uint) error {
	return database.DB.Delete(&models.Link{}, id).Error
}
