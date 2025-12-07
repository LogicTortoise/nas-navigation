package services

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"nas-navigation/database"
	"nas-navigation/models"
)

type SystemService struct{}

func NewSystemService() *SystemService {
	return &SystemService{}
}

// GetInfo returns system information
func (s *SystemService) GetInfo() (*models.SystemInfo, error) {
	hostname, _ := os.Hostname()

	info := &models.SystemInfo{
		Hostname: hostname,
		Platform: runtime.GOOS,
	}

	// Get uptime on macOS/Linux
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		cmd := exec.Command("uptime")
		output, err := cmd.Output()
		if err == nil {
			info.LoadAverage = strings.TrimSpace(string(output))
		}
	}

	return info, nil
}

// Restart performs a system restart
func (s *SystemService) Restart(target string) (*models.SystemActionResult, error) {
	result := &models.SystemActionResult{
		Action: "restart",
	}
	startTime := time.Now()

	var cmd *exec.Cmd
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Determine what to restart
	if target == "" || target == "system" {
		// System restart - use appropriate command based on OS
		switch runtime.GOOS {
		case "darwin":
			// macOS - use shutdown command (requires sudo)
			cmd = exec.CommandContext(ctx, "sudo", "shutdown", "-r", "now")
		case "linux":
			// Linux - use systemctl or shutdown
			cmd = exec.CommandContext(ctx, "sudo", "systemctl", "reboot")
		default:
			result.Success = false
			result.Error = fmt.Sprintf("unsupported platform: %s", runtime.GOOS)
			return result, nil
		}
	} else {
		// Restart specific service
		switch runtime.GOOS {
		case "darwin":
			// macOS - use launchctl
			cmd = exec.CommandContext(ctx, "sudo", "launchctl", "kickstart", "-k", fmt.Sprintf("system/%s", target))
		case "linux":
			// Linux - use systemctl
			cmd = exec.CommandContext(ctx, "sudo", "systemctl", "restart", target)
		default:
			result.Success = false
			result.Error = fmt.Sprintf("unsupported platform: %s", runtime.GOOS)
			return result, nil
		}
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result.Duration = time.Since(startTime).Milliseconds()

	if ctx.Err() == context.DeadlineExceeded {
		result.Success = false
		result.Error = "operation timed out"
		return result, nil
	}

	if err != nil {
		result.Success = false
		result.Error = stderr.String()
		if result.Error == "" {
			result.Error = err.Error()
		}
		return result, nil
	}

	result.Success = true
	if target == "" || target == "system" {
		result.Message = "System restart initiated"
	} else {
		result.Message = fmt.Sprintf("Service '%s' restarted successfully", target)
	}

	return result, nil
}

// RestartLink restarts a specific link by its ID
func (s *SystemService) RestartLink(id string) (*models.SystemActionResult, error) {
	var link models.Link
	if err := database.DB.First(&link, id).Error; err != nil {
		return nil, fmt.Errorf("link not found: %w", err)
	}

	if link.RestartScript == "" {
		return &models.SystemActionResult{
			Action:  "restart",
			Success: false,
			Error:   "该链接未配置重启脚本",
		}, nil
	}

	result := s.executeScript(link.RestartScript, 30)
	result.Action = "restart"
	result.Message = fmt.Sprintf("服务 '%s' 重启完成", link.Name)
	return result, nil
}

// RestartAllServices restarts all services that have restart scripts configured
func (s *SystemService) RestartAllServices() (*models.BatchRestartResult, error) {
	// Get all links with restart_script configured
	var links []models.Link
	err := database.DB.Where("restart_script IS NOT NULL AND restart_script != ''").Find(&links).Error
	if err != nil {
		return nil, err
	}

	result := &models.BatchRestartResult{
		TotalCount:   len(links),
		Results:      make([]models.ServiceRestartResult, 0),
		SuccessCount: 0,
		FailedCount:  0,
	}

	if len(links) == 0 {
		result.Message = "没有配置重启脚本的服务"
		return result, nil
	}

	// Execute scripts concurrently with a limit
	var wg sync.WaitGroup
	resultChan := make(chan models.ServiceRestartResult, len(links))
	semaphore := make(chan struct{}, 3) // Limit concurrent executions to 3

	for _, link := range links {
		wg.Add(1)
		go func(l models.Link) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			serviceResult := models.ServiceRestartResult{
				ServiceName: l.Name,
				ServiceID:   l.ID,
				ScriptName:  l.RestartScript,
			}

			startTime := time.Now()
			execResult := s.executeScript(l.RestartScript, 30)
			serviceResult.Duration = time.Since(startTime).Milliseconds()
			serviceResult.Success = execResult.Success
			serviceResult.Output = execResult.Message
			serviceResult.Error = execResult.Error

			resultChan <- serviceResult
		}(link)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for r := range resultChan {
		result.Results = append(result.Results, r)
		if r.Success {
			result.SuccessCount++
		} else {
			result.FailedCount++
		}
	}

	if result.FailedCount == 0 {
		result.Message = fmt.Sprintf("成功重启 %d 个服务", result.SuccessCount)
	} else {
		result.Message = fmt.Sprintf("重启完成: %d 成功, %d 失败", result.SuccessCount, result.FailedCount)
	}

	return result, nil
}

// executeScript executes a single script
func (s *SystemService) executeScript(command string, timeout int) *models.SystemActionResult {
	result := &models.SystemActionResult{
		Action: "execute",
	}

	if timeout <= 0 {
		timeout = 30
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		result.Success = false
		result.Error = "执行超时"
		return result
	}

	if err != nil {
		result.Success = false
		result.Error = stderr.String()
		if result.Error == "" {
			result.Error = err.Error()
		}
		result.Message = stdout.String()
		return result
	}

	result.Success = true
	result.Message = stdout.String()
	return result
}

// ExecuteCommand executes a custom command (for quick actions)
func (s *SystemService) ExecuteCommand(command string, timeout int) (*models.SystemActionResult, error) {
	result := &models.SystemActionResult{
		Action: "execute",
	}
	startTime := time.Now()

	if timeout <= 0 {
		timeout = 30
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result.Duration = time.Since(startTime).Milliseconds()

	if ctx.Err() == context.DeadlineExceeded {
		result.Success = false
		result.Error = "execution timed out"
		return result, nil
	}

	if err != nil {
		result.Success = false
		result.Error = stderr.String()
		if result.Error == "" {
			result.Error = err.Error()
		}
		result.Message = stdout.String()
		return result, nil
	}

	result.Success = true
	result.Message = stdout.String()
	return result, nil
}
