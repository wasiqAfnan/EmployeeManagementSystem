package routes

import (
	"net/http"

	"EMS/internal/handler"
)

func EmployeeRoutes(mux *http.ServeMux, h *handler.EmployeeHandler) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Employee Management System API is running"))
	})
	mux.HandleFunc("GET /api/employees", h.GetEmployees)
	mux.HandleFunc("GET /api/employee/{empId}", h.GetEmployee)
	mux.HandleFunc("GET /api/employees/search", h.SearchEmployees)
	mux.HandleFunc("POST /api/employee", h.CreateEmployee)
	mux.HandleFunc("PATCH /api/employee/{empId}", h.UpdateEmployee)
	mux.HandleFunc("DELETE /api/employee/{empId}", h.DeleteEmployee)
}

func UserRoutes(mux *http.ServeMux, h *handler.UserHandler) {
	mux.HandleFunc("GET /api/me", h.Me)
	mux.HandleFunc("GET /api/users", h.GetUsers)
}
