package handlers

import (
	"net/http"

	"nas-navigation/services"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	service *services.SettingService
}

func NewSettingHandler() *SettingHandler {
	return &SettingHandler{
		service: services.NewSettingService(),
	}
}

func (h *SettingHandler) GetAll(c *gin.Context) {
	settings, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (h *SettingHandler) Update(c *gin.Context) {
	var settings map[string]string
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetBatch(settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "设置已更新"})
}
