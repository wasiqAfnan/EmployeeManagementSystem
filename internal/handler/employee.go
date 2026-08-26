package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"EMS/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// EmployeeHandler holds a direct reference to the MongoDB collection.
type EmployeeHandler struct {
	collection *mongo.Collection
}

// NewEmployeeHandler creates a new EmployeeHandler with the given collection.
func NewEmployeeHandler(collection *mongo.Collection) *EmployeeHandler {
	return &EmployeeHandler{collection: collection}
}

// APIResponse is the standard JSON envelope sent back on every request.
type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// sendJSON is a helper fucntion that encodes the response.
func sendJSON(w http.ResponseWriter, statusCode int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

// GET /employees
func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Fetch all documents from the collection
	result, err := h.collection.Find(ctx, bson.M{})
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch employees: " + err.Error(),
			Data:    nil,
		})
		return
	}
	defer result.Close(ctx)

	// Decode the cursor results into a slice
	var employees []model.Employee
	if err = result.All(ctx, &employees); err != nil {
		sendJSON(w, http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to decode employees: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Return empty array instead of null when there are no records
	if employees == nil {
		employees = []model.Employee{}
	}

	sendJSON(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Employee records fetched successfully",
		Data:    employees,
	})
}
