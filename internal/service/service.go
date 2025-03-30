package service

import (
	"errors"

	"github.com/DblMOKRQ/introductory-practice/internal/entity"
	"github.com/DblMOKRQ/introductory-practice/internal/repository"
	"github.com/DblMOKRQ/introductory-practice/pkg/logger"
	"go.uber.org/zap"
)

type Service struct {
	repo   *repository.Repository
	logger *logger.Logger
}

func NewService(repo *repository.Repository, logger *logger.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) AddVehicle(v *entity.Vehicle) error {
	err := s.repo.AddVehicle(v)
	if err != nil {
		// s.logger.Error("vehicle already exists",zap.Error(err))
		return errors.New("vehicle already exists")
	}

	return nil
}

func (s *Service) GetVehicle(vin string) (*entity.Vehicle, error) {
	vehicle, err := s.repo.GetVehicle(vin)
	if err != nil {
		s.logger.Error("vehicle not found", zap.Error(err))
		return nil, errors.New("vehicle not found")
	}
	return vehicle, nil

}

func (s *Service) UpdateVehicle(v *entity.Vehicle) error {
	err := s.repo.UpdateVehicle(v)
	if err != nil {
		s.logger.Error("vehicle not found", zap.Error(err))
		return errors.New("vehicle not found")
	}
	return nil
}

func (s *Service) DeleteVehicle(vin string) error {
	err := s.repo.DeleteVehicle(vin)
	if err != nil {
		s.logger.Error("vehicle not found", zap.Error(err))
		return errors.New("vehicle not found")
	}
	return nil
}

func (s *Service) GetAllVehicles() ([]*entity.Vehicle, error) {
	vehicles, err := s.repo.GetAllVehicles()
	if err != nil {
		s.logger.Error("vehicles not found", zap.Error(err))
		return nil, errors.New("vehicles not found")
	}
	return vehicles, nil
}

func (s *Service) RentVehicle(vin string, user *entity.User) error {
	err := s.repo.RentVehicle(vin, user)
	if err != nil {
		s.logger.Error("vehicle not found", zap.String("vin", vin), zap.Error(err))
		return errors.New("vehicle not found")
	}
	return nil
}

func (s *Service) UpdateStatus(vin string, status string) error {
	err := s.repo.UpdateStatus(vin, status)
	if err != nil {
		s.logger.Error("vehicle not found", zap.String("vin", vin), zap.Error(err))
		return errors.New("vehicle not found")
	}
	return nil
}
