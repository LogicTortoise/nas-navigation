package database

import (
	"log"
	"os"
	"path/filepath"

	"nas-navigation/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(dbPath string) error {
	// 确保目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// 自动迁移
	if err := DB.AutoMigrate(
		&models.Category{},
		&models.Link{},
		&models.Script{},
		&models.Setting{},
	); err != nil {
		return err
	}

	// 初始化默认数据
	if err := initDefaultData(); err != nil {
		log.Printf("Warning: failed to init default data: %v", err)
	}

	return nil
}

func initDefaultData() error {
	// 检查是否已有分类
	var count int64
	DB.Model(&models.Category{}).Count(&count)
	if count > 0 {
		return nil
	}

	// 创建默认分类
	categories := []models.Category{
		{Name: "常用访问", Icon: "star", SortOrder: 1},
		{Name: "工作效率", Icon: "briefcase", SortOrder: 2},
		{Name: "设计灵感", Icon: "palette", SortOrder: 3},
	}

	for _, cat := range categories {
		if err := DB.Create(&cat).Error; err != nil {
			return err
		}
	}

	// 创建默认设置
	settings := []models.Setting{
		{Key: models.SettingUserName, Value: "用户"},
		{Key: models.SettingGreeting, Value: "准备好开始高效的一天了吗?"},
		{Key: models.SettingTheme, Value: "light"},
		{Key: models.SettingLanguage, Value: "zh-CN"},
	}

	for _, setting := range settings {
		DB.FirstOrCreate(&setting, models.Setting{Key: setting.Key})
	}

	return nil
}

func GetDB() *gorm.DB {
	return DB
}
