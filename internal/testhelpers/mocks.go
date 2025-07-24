package testhelpers

import (
	"computerInventory/internal/domain"
	"errors"
	"fmt"

	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	store map[string]domain.Computer
}

func NewMockRepo() *MockRepo {
	return &MockRepo{store: make(map[string]domain.Computer)}
}

func (m *MockRepo) Create(c *domain.Computer) error {
	m.store[c.MACAddress] = *c
	fmt.Println("Create for emp:" + c.EmployeeAbbreviation)
	return nil
}
func (m *MockRepo) Get(mac string) (*domain.Computer, error) {
	if c, ok := m.store[mac]; ok {
		return &c, nil
	}
	return nil, errors.New("not found")
}
func (m *MockRepo) Update(c *domain.Computer) error {
	m.store[c.MACAddress] = *c
	return nil
}
func (m *MockRepo) Delete(mac string) error {
	delete(m.store, mac)
	return nil
}
func (m *MockRepo) GetAll() ([]domain.Computer, error) {
	var all []domain.Computer
	for _, c := range m.store {
		all = append(all, c)
	}
	return all, nil
}
func (m *MockRepo) GetByEmployee(abbr string) ([]domain.Computer, error) {
	var result []domain.Computer

	for _, c := range m.store {

		if c.EmployeeAbbreviation == abbr {
			result = append(result, c)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("not found")
	}

	fmt.Printf("GetByEmployee:%v %v\n", abbr, len(result))

	return result, nil
}

type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) SendWarning(abbr string, message string) error {
	m.Called(abbr, message)
	return nil
}
