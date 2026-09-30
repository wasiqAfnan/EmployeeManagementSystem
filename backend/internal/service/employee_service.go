package service

import (
	"context"
	"errors"
	"strings"

	"EMS/internal/model"
	"EMS/internal/repository"
	"EMS/internal/utils"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrEmployeeNotFound    = errors.New("employee not found")
	ErrInvalidEmpID        = errors.New("empId is required")
	ErrDuplicateEmpID      = errors.New("empId already exists")
	ErrSearchQueryRequired = errors.New("search query is required")
)

type EmployeeService struct {
	repo *repository.EmployeeRepository
}

func NewEmployeeService(repo *repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

// GetAllEmployees handles business logic for fetching all employees.
func (s *EmployeeService) GetAllEmployees(ctx context.Context, sub string) ([]model.Employee, error) {
	return s.repo.GetAll(ctx, sub)
}

// GetEmployeeByEmpID handles validation and retrieval logic for a single employee by empId.
func (s *EmployeeService) GetEmployeeByEmpID(ctx context.Context, sub string, empID string) (*model.Employee, error) {
	cleanID := strings.TrimSpace(empID)
	if cleanID == "" {
		return nil, ErrInvalidEmpID
	}

	employee, err := s.repo.GetByEmpID(ctx, sub, cleanID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	return employee, nil
}

// CreateEmployee handles validation, duplicate check, and creation logic for a new employee.
func (s *EmployeeService) CreateEmployee(ctx context.Context, emp model.Employee) (*model.Employee, error) {
	if err := utils.ValidateEmployee(emp); err != nil {
		return nil, err
	}

	// Check if employee with given empId already exists for this user
	_, err := s.repo.GetByEmpID(ctx, emp.CreatedBy, emp.EmpID)
	if err == nil {
		return nil, ErrDuplicateEmpID
	}

	return s.repo.Create(ctx, &emp)
}

// UpdateEmployee handles validation and update logic for an existing employee.
func (s *EmployeeService) UpdateEmployee(ctx context.Context, sub string, empID string, emp model.Employee) (*model.Employee, error) {
	cleanID := strings.TrimSpace(empID)
	if cleanID == "" {
		return nil, ErrInvalidEmpID
	}

	if err := utils.ValidateEmployee(emp); err != nil {
		return nil, err
	}

	updatedEmp, err := s.repo.Update(ctx, sub, cleanID, emp)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	return updatedEmp, nil
}

// DeleteEmployee handles deletion logic for an employee by empId.
func (s *EmployeeService) DeleteEmployee(ctx context.Context, sub string, empID string) (*model.Employee, error) {
	cleanID := strings.TrimSpace(empID)
	if cleanID == "" {
		return nil, ErrInvalidEmpID
	}

	deletedEmp, err := s.repo.Delete(ctx, sub, cleanID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	return deletedEmp, nil
}

// SearchEmployees handles search validation and query logic across employees.
func (s *EmployeeService) SearchEmployees(ctx context.Context, sub string, query string) ([]model.Employee, error) {
	cleanQuery := strings.TrimSpace(query)
	if cleanQuery == "" {
		return nil, ErrSearchQueryRequired
	}

	return s.repo.Search(ctx, sub, cleanQuery)
}
