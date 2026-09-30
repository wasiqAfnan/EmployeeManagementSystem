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
	repo     *repository.EmployeeRepository
	userRepo *repository.UserRepository
}

func NewEmployeeService(repo *repository.EmployeeRepository, userRepo *repository.UserRepository) *EmployeeService {
	return &EmployeeService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *EmployeeService) populateCreatorEmails(ctx context.Context, employees []model.Employee) []model.Employee {
	if s.userRepo == nil || len(employees) == 0 {
		return employees
	}

	users, err := s.userRepo.GetAll(ctx)
	if err != nil || len(users) == 0 {
		return employees
	}

	userMap := make(map[string]string)
	for _, u := range users {
		userMap[u.Sub] = u.Email
	}

	for i := range employees {
		if email, found := userMap[employees[i].CreatedBy]; found {
			employees[i].CreatedByEmail = email
		}
	}

	return employees
}

// GetAllEmployees handles business logic for fetching all employees.
func (s *EmployeeService) GetAllEmployees(ctx context.Context, sub string, role string) ([]model.Employee, error) {
	employees, err := s.repo.GetAll(ctx, sub, role)
	if err != nil {
		return nil, err
	}

	if role == "admin" {
		employees = s.populateCreatorEmails(ctx, employees)
	}

	return employees, nil
}

// GetEmployeeByEmpID handles validation and retrieval logic for a single employee by empId.
func (s *EmployeeService) GetEmployeeByEmpID(ctx context.Context, sub string, role string, empID string) (*model.Employee, error) {
	cleanID := strings.TrimSpace(empID)
	if cleanID == "" {
		return nil, ErrInvalidEmpID
	}

	employee, err := s.repo.GetByEmpID(ctx, sub, role, cleanID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	if role == "admin" && employee != nil {
		employees := s.populateCreatorEmails(ctx, []model.Employee{*employee})
		employee = &employees[0]
	}

	return employee, nil
}

// CreateEmployee handles validation, duplicate check, and creation logic for a new employee.
func (s *EmployeeService) CreateEmployee(ctx context.Context, emp model.Employee) (*model.Employee, error) {
	if err := utils.ValidateEmployee(emp); err != nil {
		return nil, err
	}

	// Check if employee with given empId already exists for this user
	_, err := s.repo.GetByEmpID(ctx, emp.CreatedBy, "user", emp.EmpID)

	// If no error, it means the employee already exists
	if err == nil {
		return nil, ErrDuplicateEmpID
	}

	return s.repo.Create(ctx, &emp)
}

// UpdateEmployee handles validation and update logic for an existing employee.
func (s *EmployeeService) UpdateEmployee(ctx context.Context, sub string, role string, empID string, emp model.Employee) (*model.Employee, error) {
	cleanID := strings.TrimSpace(empID)
	if cleanID == "" {
		return nil, ErrInvalidEmpID
	}

	if err := utils.ValidateEmployee(emp); err != nil {
		return nil, err
	}

	updatedEmp, err := s.repo.Update(ctx, sub, role, cleanID, emp)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	if role == "admin" && updatedEmp != nil {
		employees := s.populateCreatorEmails(ctx, []model.Employee{*updatedEmp})
		updatedEmp = &employees[0]
	}

	return updatedEmp, nil
}

// DeleteEmployee handles deletion logic for an employee by empId.
func (s *EmployeeService) DeleteEmployee(ctx context.Context, sub string, role string, empID string) (*model.Employee, error) {
	cleanID := strings.TrimSpace(empID)
	if cleanID == "" {
		return nil, ErrInvalidEmpID
	}

	deletedEmp, err := s.repo.Delete(ctx, sub, role, cleanID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	return deletedEmp, nil
}

// SearchEmployees handles search validation and query logic across employees.
func (s *EmployeeService) SearchEmployees(ctx context.Context, sub string, role string, query string) ([]model.Employee, error) {
	cleanQuery := strings.TrimSpace(query)
	if cleanQuery == "" {
		return nil, ErrSearchQueryRequired
	}

	employees, err := s.repo.Search(ctx, sub, role, cleanQuery)
	if err != nil {
		return nil, err
	}

	if role == "admin" {
		employees = s.populateCreatorEmails(ctx, employees)
	}

	return employees, nil
}
