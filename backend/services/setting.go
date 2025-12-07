package services

import (
	"nas-navigation/database"
	"nas-navigation/models"
)

type SettingService struct{}

func NewSettingService() *SettingService {
	return &SettingService{}
}

func (s *SettingService) GetAll() (map[string]string, error) {
	var settings []models.Setting
	err := database.DB.Find(&settings).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}
	return result, nil
}

func (s *SettingService) Get(key string) (string, error) {
	var setting models.Setting
	err := database.DB.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (s *SettingService) Set(key, value string) error {
	var setting models.Setting
	result := database.DB.Where("key = ?", key).First(&setting)

	if result.Error != nil {
		// 不存在则创建
		setting = models.Setting{Key: key, Value: value}
		return database.DB.Create(&setting).Error
	}

	// 存在则更新
	setting.Value = value
	return database.DB.Save(&setting).Error
}

func (s *SettingService) SetBatch(settings map[string]string) error {
	for key, value := range settings {
		if err := s.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}
