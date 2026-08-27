package routes

import (
	"net/http"

	"EMS/internal/handler"
)

func EmployeeRoutes(mux *http.ServeMux, h *handler.EmployeeHandler) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Employee Management System API is running"))
	})
	mux.HandleFunc("GET /employees", h.GetEmployees)
	mux.HandleFunc("GET /employee/{empId}", h.GetEmployee)
	mux.HandleFunc("GET /employees/search", h.SearchEmployees)
	mux.HandleFunc("POST /employee", h.CreateEmployee)
	mux.HandleFunc("PATCH /employee/{empId}", h.UpdateEmployee)
	mux.HandleFunc("DELETE /employee/{empId}", h.DeleteEmployee)
}
