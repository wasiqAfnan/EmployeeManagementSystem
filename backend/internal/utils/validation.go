package utils

import (
	"errors"
	"strings"

	"EMS/internal/model"
)

func ValidateEmployee(emp model.Employee) error {
	if strings.TrimSpace(emp.EmpID) == "" {
		return errors.New("empId is required")
	}

	if strings.TrimSpace(emp.FullName) == "" {
		return errors.New("fullName is required")
	}

	if strings.TrimSpace(emp.JobTitle) == "" {
		return errors.New("jobTitle is required")
	}

	if strings.TrimSpace(emp.Department) == "" {
		return errors.New("department is required")
	}

	if emp.Salary <= 0 {
		return errors.New("salary must be greater than 0")
	}

	return nil
}
