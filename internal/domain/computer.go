package domain

type Computer struct {
	MACAddress           string `json:"mac_address" binding:"required" gorm:"uniqueIndex;primaryKey"`
	ComputerName         string `json:"computer_name" binding:"required"`
	IPAddress            string `json:"ip_address" binding:"required"`
	EmployeeAbbreviation string `json:"employee_abbreviation,omitempty"`
	Description          string `json:"description,omitempty"`
}

type ComputerRepository interface {
	Create(computer *Computer) error
	Get(mac string) (*Computer, error)
	Update(computer *Computer) error
	Delete(mac string) error
	GetAll() ([]Computer, error)
	GetByEmployee(abbr string) ([]Computer, error)
}
