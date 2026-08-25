package routes

import (
	"net/http"

	controllers "EMS/controllers"
)

func EmployeeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("The server is up and running"))
	})
	mux.HandleFunc("GET /employees", controllers.GetEmployees)
	mux.HandleFunc("POST /employee", controllers.CreateEmployee)
	mux.HandleFunc("PUT /employee/{empId}", controllers.UpdateEmployee)
	mux.HandleFunc("DELETE /employee/{empId}", controllers.DeleteEmployee)
	// mux.HandleFunc("GET /employee/{empId}", controllers.GetEmployee)

}
