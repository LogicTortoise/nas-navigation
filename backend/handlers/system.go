package handlers

import (
	"net/http"

	"nas-navigation/models"
	"nas-navigation/services"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	service *services.SystemService
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{
		service: services.NewSystemService(),
	}
}

// GetInfo returns system information
func (h *SystemHandler) GetInfo(c *gin.Context) {
	info, err := h.service.GetInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// Restart handles system or service restart
func (h *SystemHandler) Restart(c *gin.Context) {
	var req models.SystemActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If no body, default to system restart
		req.Target = ""
	}

	// Safety check: require confirmation for system-wide restart
	if (req.Target == "" || req.Target == "system") && !req.Confirm {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "confirmation required",
			"message": "Set 'confirm: true' to proceed with system restart",
		})
		return
	}

	result, err := h.service.Restart(req.Target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// RestartAllServices restarts all services with associated scripts
func (h *SystemHandler) RestartAllServices(c *gin.Context) {
	result, err := h.service.RestartAllServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// RestartLink restarts a specific link by executing its restart_script
func (h *SystemHandler) RestartLink(c *gin.Context) {
	id := c.Param("id")
	result, err := h.service.RestartLink(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ExecuteQuickAction executes a predefined quick action
func (h *SystemHandler) ExecuteQuickAction(c *gin.Context) {
	action := c.Param("action")

	// Predefined safe quick actions
	actions := map[string]string{
		"clear-cache":    "sync && echo 'Cache cleared'",
		"disk-usage":     "df -h",
		"memory-usage":   "free -h 2>/dev/null || vm_stat",
		"network-status": "netstat -an | head -20",
		"process-list":   "ps aux | head -20",
	}

	command, exists := actions[action]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "unknown action",
			"available_actions": []string{"clear-cache", "disk-usage", "memory-usage", "network-status", "process-list"},
		})
		return
	}

	result, err := h.service.ExecuteCommand(command, 30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result.Action = action
	c.JSON(http.StatusOK, result)
}
