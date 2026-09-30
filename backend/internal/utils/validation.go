package utils

import (
	"errors"
	"regexp"
	"strings"

	"EMS/internal/model"
)

var empIDRegex = regexp.MustCompile(`^E\d{3}$`)

func ValidateEmployee(emp model.Employee) error {
	trimmedEmpID := strings.TrimSpace(emp.EmpID)
	if trimmedEmpID == "" {
		return errors.New("empId is required")
	}

	if !empIDRegex.MatchString(trimmedEmpID) {
		return errors.New("empId must start with 'E' followed by exactly 3 digits (e.g., E001)")
	}

	if strings.TrimSpace(emp.FullName) == "" {
		return errors.New("fullName is required")
	}

	email := strings.TrimSpace(emp.Email)
	if email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(email, "@") {
		return errors.New("email is invalid")
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
