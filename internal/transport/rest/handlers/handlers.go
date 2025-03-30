package handlers

import (
	"github.com/DblMOKRQ/introductory-practice/internal/entity"
	"github.com/DblMOKRQ/introductory-practice/internal/service"
	"github.com/DblMOKRQ/introductory-practice/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handlers struct {
	logger *logger.Logger
	server *service.Service
}

func NewHandlers(logger *logger.Logger, server *service.Service) *Handlers {
	return &Handlers{logger: logger, server: server}
}

func (h *Handlers) AddVehicle(c *gin.Context) {
	var v entity.Vehicle

	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Adding vehicle", zap.Any("vehicle", v))
	if err := h.server.AddVehicle(&v); err != nil {
		h.logger.Error("Error adding vehicle", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Vehicle added")
	c.JSON(200, v)

}

func (h *Handlers) GetVehicle(c *gin.Context) {
	vin := c.Param("vin")
	h.logger.Info("Getting vehicle", zap.String("vin", vin))
	v, err := h.server.GetVehicle(vin)
	if err != nil {
		h.logger.Error("Error getting vehicle", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Vehicle found", zap.Any("vehicle", v))
	c.JSON(200, v)
}

func (h *Handlers) DeleteVehicle(c *gin.Context) {
	vin := c.Param("vin")
	h.logger.Info("Deleting vehicle", zap.String("vin", vin))
	if err := h.server.DeleteVehicle(vin); err != nil {
		h.logger.Error("Error deleting vehicle", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Vehicle deleted")
	c.JSON(200, gin.H{"message": "Vehicle deleted"})
}

func (h *Handlers) GetAllVehicles(c *gin.Context) {
	h.logger.Info("Getting all vehicles")
	v, err := h.server.GetAllVehicles()
	if err != nil {
		h.logger.Error("Error getting all vehicles", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Vehicles found", zap.Any("vehicles", v))
	c.JSON(200, v)
}

func (h *Handlers) RentVehicle(c *gin.Context) {
	vin := c.Param("vin")
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Renting vehicle", zap.Any("user", user))

	if err := h.server.RentVehicle(vin, &user); err != nil {

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Vehicle was successfully booking", zap.Any("user", user))

	c.JSON(200, gin.H{"message": "Бронирование успешно"})
}

func (h *Handlers) UpdateStatusVehicle(c *gin.Context) {
	vin := c.Param("vin")
	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Updating status", zap.String("status", statusUpdate.Status), zap.String("vin", vin))

	validStatuses := map[string]bool{
		"available":         true,
		"on_route":          true,
		"under_maintenance": true,
	}

	if !validStatuses[statusUpdate.Status] {
		c.JSON(400, gin.H{"error": "Некорректный статус"})
		return
	}

	if err := h.server.UpdateStatus(vin, statusUpdate.Status); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Status updated", zap.String("status", statusUpdate.Status), zap.String("vin", vin))

	c.JSON(200, gin.H{"message": "Статус обновлен"})
}
func (h *Handlers) CheckAdmin(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Status(200)
	c.JSON(200, gin.H{"message": "OK"})
}
func (h *Handlers) AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const adminUser = "admin"
		const adminPass = "1"

		user, pass, ok := c.Request.BasicAuth()
		if !ok || user != adminUser || pass != adminPass {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(401, gin.H{"error": "Требуется аутентификация"})
			return
		}
		c.Next()
	}
}
