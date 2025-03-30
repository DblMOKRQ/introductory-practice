package rout

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type handlers interface {
	AddVehicle(c *gin.Context)
	GetVehicle(c *gin.Context)
	// UpdateVehicle(c *gin.Context)
	DeleteVehicle(c *gin.Context)
	GetAllVehicles(c *gin.Context)
	RentVehicle(c *gin.Context)
	UpdateStatusVehicle(c *gin.Context)
	CheckAdmin(c *gin.Context)
	AdminAuthMiddleware() gin.HandlerFunc
}
type Rout struct {
	router *gin.Engine
	h      handlers
}

func NewRout(hand handlers) *Rout {
	r := gin.Default()
	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization", // Добавляем заголовок авторизации
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return &Rout{router: r, h: hand}
}

func (r *Rout) Run(host string, port string) error {

	r.router.POST("/add", r.h.AddVehicle)

	r.router.GET("/get/:vin", r.h.GetVehicle)

	// r.router.PUT("/update", func(c *gin.Context) {
	// 	var v entity.Vehicle

	// 	if err := c.ShouldBindJSON(&v); err != nil {
	// 		c.JSON(400, gin.H{"error": err.Error()})
	// 		r.logger.Error("Error binding JSON", zap.Error(err))
	// 		return
	// 	}
	// 	r.logger.Info("Updating vehicle", zap.Any("vehicle", v))
	// 	if err := r.server.UpdateVehicle(&v); err != nil {
	// 		r.logger.Error("Error updating vehicle", zap.Error(err))
	// 		c.JSON(400, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	c.JSON(200, v)
	// })

	r.router.DELETE("/delete/:vin", r.h.AdminAuthMiddleware(), r.h.DeleteVehicle)

	r.router.GET("/getall", r.h.GetAllVehicles)

	r.router.POST("/rent/:vin", r.h.RentVehicle)

	r.router.PUT("/updateStatus/:vin", r.h.AdminAuthMiddleware(), r.h.UpdateStatusVehicle)

	r.router.GET("/admin/check", r.h.AdminAuthMiddleware(), r.h.CheckAdmin)

	return r.router.Run(":8080")
}
