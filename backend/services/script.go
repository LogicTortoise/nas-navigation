package services

import (
	"bytes"
	"context"
	"os/exec"
	"time"

	"nas-navigation/database"
	"nas-navigation/models"
)

type ScriptService struct{}

func NewScriptService() *ScriptService {
	return &ScriptService{}
}

func (s *ScriptService) GetAll() ([]models.Script, error) {
	var scripts []models.Script
	err := database.DB.Order("name ASC").Find(&scripts).Error
	return scripts, err
}

func (s *ScriptService) GetByID(id uint) (*models.Script, error) {
	var script models.Script
	err := database.DB.First(&script, id).Error
	if err != nil {
		return nil, err
	}
	return &script, nil
}

func (s *ScriptService) Create(req *models.ScriptRequest) (*models.Script, error) {
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	script := &models.Script{
		Name:        req.Name,
		Command:     req.Command,
		Description: req.Description,
		Timeout:     timeout,
	}
	err := database.DB.Create(script).Error
	return script, err
}

func (s *ScriptService) Update(id uint, req *models.ScriptRequest) (*models.Script, error) {
	var script models.Script
	if err := database.DB.First(&script, id).Error; err != nil {
		return nil, err
	}

	script.Name = req.Name
	script.Command = req.Command
	script.Description = req.Description
	if req.Timeout > 0 {
		script.Timeout = req.Timeout
	}

	err := database.DB.Save(&script).Error
	return &script, err
}

func (s *ScriptService) Delete(id uint) error {
	// 清除关联的链接中的脚本引用
	database.DB.Model(&models.Link{}).Where("script_id = ?", id).Update("script_id", nil)
	return database.DB.Delete(&models.Script{}, id).Error
}

func (s *ScriptService) Execute(id uint) (*models.ScriptExecuteResult, error) {
	script, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	result := &models.ScriptExecuteResult{}
	startTime := time.Now()

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(script.Timeout)*time.Second)
	defer cancel()

	// 执行命令
	cmd := exec.CommandContext(ctx, "sh", "-c", script.Command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	result.Duration = time.Since(startTime).Milliseconds()
	result.Output = stdout.String()

	if ctx.Err() == context.DeadlineExceeded {
		result.Success = false
		result.Error = "执行超时"
		result.ExitCode = -1
		return result, nil
	}

	if err != nil {
		result.Success = false
		result.Error = stderr.String()
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		return result, nil
	}

	result.Success = true
	result.ExitCode = 0
	return result, nil
}
