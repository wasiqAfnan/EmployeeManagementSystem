package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"EMS/internal/middleware"
	"EMS/internal/model"
	"EMS/internal/service"
	"EMS/internal/utils"
)

// EmployeeHandler handles employee HTTP requests by delegating to EmployeeService.
type EmployeeHandler struct {
	service *service.EmployeeService
}

// NewEmployeeHandler creates a new EmployeeHandler.
func NewEmployeeHandler(service *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service: service,
	}
}

// GET /employees
func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	sub, _ := r.Context().Value(middleware.UserContextKey).(string)
	role, _ := r.Context().Value(middleware.UserRoleContextKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	employees, err := h.service.GetAllEmployees(ctx, sub, role)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to fetch employees: " + err.Error(),
			Data:    nil,
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Employee records fetched successfully",
		Data:    employees,
	})
}

// GET /employee/{empId}
func (h *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	empId := r.PathValue("empId")
	sub, _ := r.Context().Value(middleware.UserContextKey).(string)
	role, _ := r.Context().Value(middleware.UserRoleContextKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	employee, err := h.service.GetEmployeeByEmpID(ctx, sub, role, empId)
	if err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			utils.SendJSON(w, http.StatusNotFound, utils.APIResponse{
				Status:  "error",
				Message: "Employee not found",
				Data:    nil,
			})
			return
		}
		if errors.Is(err, service.ErrInvalidEmpID) {
			utils.SendJSON(w, http.StatusBadRequest, utils.APIResponse{
				Status:  "error",
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		utils.SendJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to fetch employee: " + err.Error(),
			Data:    nil,
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Employee record fetched successfully",
		Data:    employee,
	})
}

// POST /employee
func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp model.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid JSON body",
			Data:    nil,
		})
		return
	}

	// Extract sub from context and set CreatedBy in the payload
	if sub, ok := r.Context().Value(middleware.UserContextKey).(string); ok {
		emp.CreatedBy = sub
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	createdEmp, err := h.service.CreateEmployee(ctx, emp)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateEmpID) {
			utils.SendJSON(w, http.StatusBadRequest, utils.APIResponse{
				Status:  "error",
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		// Validation error or other business logic error
		utils.SendJSON(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	utils.SendJSON(w, http.StatusCreated, utils.APIResponse{
		Status:  "success",
		Message: "Employee record created successfully",
		Data:    createdEmp,
	})
}

// PUT /employee/{empId}
func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	empId := r.PathValue("empId")
	sub, _ := r.Context().Value(middleware.UserContextKey).(string)
	role, _ := r.Context().Value(middleware.UserRoleContextKey).(string)

	var emp model.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid JSON body",
			Data:    nil,
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	updatedEmp, err := h.service.UpdateEmployee(ctx, sub, role, empId, emp)
	if err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			utils.SendJSON(w, http.StatusNotFound, utils.APIResponse{
				Status:  "error",
				Message: "Employee not found",
				Data:    nil,
			})
			return
		}
		if errors.Is(err, service.ErrInvalidEmpID) {
			utils.SendJSON(w, http.StatusBadRequest, utils.APIResponse{
				Status:  "error",
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		utils.SendJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to update employee: " + err.Error(),
			Data:    nil,
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Employee record updated successfully",
		Data:    updatedEmp,
	})
}

// DELETE /employee/{empId}
func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	empId := r.PathValue("empId")
	sub, _ := r.Context().Value(middleware.UserContextKey).(string)
	role, _ := r.Context().Value(middleware.UserRoleContextKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	_, err := h.service.DeleteEmployee(ctx, sub, role, empId)
	if err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			utils.SendJSON(w, http.StatusNotFound, utils.APIResponse{
				Status:  "error",
				Message: "Employee not found",
				Data:    nil,
			})
			return
		}
		if errors.Is(err, service.ErrInvalidEmpID) {
			utils.SendJSON(w, http.StatusBadRequest, utils.APIResponse{
				Status:  "error",
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		utils.SendJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to delete employee: " + err.Error(),
			Data:    nil,
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Employee record deleted successfully",
		Data:    nil,
	})
}

// GET /employees/search?q={searchQuery}
func (h *EmployeeHandler) SearchEmployees(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	sub, _ := r.Context().Value(middleware.UserContextKey).(string)
	role, _ := r.Context().Value(middleware.UserRoleContextKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	employees, err := h.service.SearchEmployees(ctx, sub, role, query)
	if err != nil {
		if errors.Is(err, service.ErrSearchQueryRequired) {
			utils.SendJSON(w, http.StatusBadRequest, utils.APIResponse{
				Status:  "error",
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		utils.SendJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to search employees: " + err.Error(),
			Data:    nil,
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Employee search completed successfully",
		Data:    employees,
	})
}
