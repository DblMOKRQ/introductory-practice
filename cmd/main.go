package main

import (
	"fmt"

	"github.com/DblMOKRQ/introductory-practice/internal/config"
	"github.com/DblMOKRQ/introductory-practice/internal/repository"
	"github.com/DblMOKRQ/introductory-practice/internal/service"
	"github.com/DblMOKRQ/introductory-practice/internal/storage/users"
	veh "github.com/DblMOKRQ/introductory-practice/internal/storage/vehicle"
	rout "github.com/DblMOKRQ/introductory-practice/internal/transport/rest"
	"github.com/DblMOKRQ/introductory-practice/internal/transport/rest/handlers"
	"github.com/DblMOKRQ/introductory-practice/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()
	log := logger.NewLogger()
	storageVehicles, err := veh.NewStorage(cfg.User, cfg.Password, cfg.DBName, cfg.Sslmode)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	storageUsers, err := users.NewStorage(cfg.User, cfg.Password, cfg.DBName, cfg.Sslmode)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	repo := repository.NewRepository(storageVehicles, storageUsers)

	srv := service.NewService(repo, log)

	hand := handlers.NewHandlers(log, srv)

	router := rout.NewRout(hand)

	log.Info("Server started", zap.String("addres", fmt.Sprint(cfg.Host, ":", cfg.Port)))
	if err := router.Run(cfg.Host, cfg.Port); err != nil {
		fmt.Println(err)
		log.Fatal("failed to start server", zap.Error(err))
	}

}
