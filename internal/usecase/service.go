package usecase

import (
	"fmt"
	"log"

	"computerInventory/internal/domain"
	"computerInventory/internal/notifier"
)

type Service struct {
	repo     domain.ComputerRepository
	notifier notifier.NotifierInterface
}

func NewService(repo domain.ComputerRepository, notifier notifier.NotifierInterface) *Service {
	return &Service{repo: repo, notifier: notifier}
}

func (s *Service) maybeNotifyEmployeeHasTooManyComputers(abbr string) {
	computers, err := s.repo.GetByEmployee(abbr)
	if err != nil {
		log.Printf("Failed to fetch computers for %s: %v", abbr, err)
		return
	}
	if len(computers) >= 3 {
		msg := fmt.Sprintf("Employee %s has %d computers", abbr, len(computers))
		if err := s.notifier.SendWarning(abbr, msg); err != nil {
			log.Printf("Warning notification failed: %v", err)
		}
	}
}

func (s *Service) Update(c *domain.Computer) error {
	if err := s.repo.Update(c); err != nil {
		return err
	}
	if c.EmployeeAbbreviation != "" {
		s.maybeNotifyEmployeeHasTooManyComputers(c.EmployeeAbbreviation)
	}
	return nil
}

func (s *Service) AddComputer(c *domain.Computer) error {
	if err := s.repo.Create(c); err != nil {
		return err
	}
	if c.EmployeeAbbreviation != "" {
		s.maybeNotifyEmployeeHasTooManyComputers(c.EmployeeAbbreviation)
	}
	return nil
}

func (s *Service) GetAll() ([]domain.Computer, error) {
	return s.repo.GetAll()
}

func (s *Service) Get(mac string) (*domain.Computer, error) {
	return s.repo.Get(mac)
}

func (s *Service) Delete(mac string) error {
	return s.repo.Delete(mac)
}

func (s *Service) GetByEmployee(abbr string) ([]domain.Computer, error) {
	return s.repo.GetByEmployee(abbr)
}
