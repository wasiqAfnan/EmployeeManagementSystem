package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"EMS/internal/model"
	"EMS/internal/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

// GET /employee/{empId}
func (h *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	empId := r.PathValue("empId")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Find one document where empId matches
	var employee model.Employee
	err := h.collection.FindOne(ctx, bson.M{"empId": empId}).Decode(&employee)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			sendJSON(w, http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: "Employee not found",
				Data:    nil,
			})
			return
		}
		sendJSON(w, http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch employee: " + err.Error(),
			Data:    nil,
		})
		return
	}

	sendJSON(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Employee record fetched successfully",
		Data:    employee,
	})
}

// POST /employee
func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON body
	var emp model.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid JSON body",
			Data:    nil,
		})
		return
	}

	// Validate required fields
	if err := utils.ValidateEmployee(emp); err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	// Keep the 10-second context
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Check duplicate empId
	var existingEmployee model.Employee

	err := h.collection.FindOne(
		ctx,
		bson.M{"empId": emp.EmpID},
	).Decode(&existingEmployee)

	if err == nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "empId already exists",
			Data:    nil,
		})
		return
	}

	// Insert employee
	var result *mongo.InsertOneResult
	result, err = h.collection.InsertOne(ctx, &emp)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to create employee: " + err.Error(),
			Data:    nil,
		})
		return
	}

	emp.ID = result.InsertedID.(bson.ObjectID)
	// Send back the created record
	sendJSON(w, http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "Employee record created successfully",
		Data:    emp,
	})
}

// PUT /employee/{empId}
func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	empId := r.PathValue("empId")

	// Parse the JSON body
	var emp model.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid JSON body",
			Data:    nil,
		})
		return
	}

	// Validate required fields
	if err := utils.ValidateEmployee(emp); err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Build the filter and update document
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	filter := bson.M{"empId": empId}
	update := bson.M{
		"$set": bson.M{
			"fullName":   emp.FullName,
			"jobTitle":   emp.JobTitle,
			"department": emp.Department,
			"salary":     emp.Salary,
		},
	}

	// Update and return the updated document
	var updatedEmp model.Employee

	err := h.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedEmp)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			sendJSON(w, http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: "Employee not found",
				Data:    nil,
			})
			return
		}

		sendJSON(w, http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to update employee: " + err.Error(),
			Data:    nil,
		})
		return
	}

	sendJSON(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Employee record updated successfully",
		Data:    updatedEmp,
	})
}

// DELETE /employee/{empId}
func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	empId := r.PathValue("empId")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Delete and return the deleted document
	var deletedEmp model.Employee

	err := h.collection.FindOneAndDelete(
		ctx,
		bson.M{"empId": empId},
	).Decode(&deletedEmp)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			sendJSON(w, http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: "Employee not found",
				Data:    nil,
			})
			return
		}

		sendJSON(w, http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to delete employee: " + err.Error(),
			Data:    nil,
		})
		return
	}

	sendJSON(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Employee record deleted successfully",
		Data:    nil,
	})
}

// GET /employees/search?q={searchQuery}
func (h *EmployeeHandler) SearchEmployees(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	search := strings.TrimSpace(r.URL.Query().Get("q"))

	if search == "" {
		sendJSON(w, http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Search query is required",
			Data:    nil,
		})
		return
	}

	// Search across multiple employee fields
	filter := bson.M{
		"$or": []bson.M{
			{"empId": bson.M{
				"$regex":   search,
				"$options": "i",
			}},
			{"fullName": bson.M{
				"$regex":   search,
				"$options": "i",
			}},
			{"jobTitle": bson.M{
				"$regex":   search,
				"$options": "i",
			}},
			{"department": bson.M{
				"$regex":   search,
				"$options": "i",
			}},
		},
	}

	var employees []model.Employee

	cursor, err := h.collection.Find(ctx, filter)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to search employees: " + err.Error(),
			Data:    nil,
		})
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &employees); err != nil {
		sendJSON(w, http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to decode employees: " + err.Error(),
			Data:    nil,
		})
		return
	}

	if employees == nil {
		employees = []model.Employee{}
	}

	sendJSON(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Employee search completed successfully",
		Data:    employees,
	})
}
