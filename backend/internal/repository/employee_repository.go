package repository

import (
	"context"
	"time"

	"EMS/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type EmployeeRepository struct {
	collection *mongo.Collection
}

func NewEmployeeRepository(collection *mongo.Collection) *EmployeeRepository {
	return &EmployeeRepository{collection: collection}
}

// GetAll fetches employee records from MongoDB (all employees if admin, tenant-only if user).
func (r *EmployeeRepository) GetAll(ctx context.Context, sub string, role string) ([]model.Employee, error) {
	findOptions := options.Find().SetSort(bson.D{{Key: "empId", Value: 1}})

	filter := bson.M{}
	if role != "admin" {
		filter["createdBy"] = sub
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var employees []model.Employee
	if err := cursor.All(ctx, &employees); err != nil {
		return nil, err
	}

	if employees == nil {
		employees = []model.Employee{}
	}

	return employees, nil
}

// GetByEmpID fetches a single employee record by empId from MongoDB.
func (r *EmployeeRepository) GetByEmpID(ctx context.Context, sub string, role string, empID string) (*model.Employee, error) {
	filter := bson.M{"empId": empID}
	if role != "admin" {
		filter["createdBy"] = sub
	}

	var employee model.Employee
	err := r.collection.FindOne(ctx, filter).Decode(&employee)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// Create inserts a new employee record into MongoDB.
func (r *EmployeeRepository) Create(ctx context.Context, emp *model.Employee) (*model.Employee, error) {
	now := time.Now()
	emp.CreatedAt = now
	emp.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, emp)
	if err != nil {
		return nil, err
	}

	// Extract the inserted _id from the result
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		emp.ID = oid
	}
	return emp, nil
}

// Update updates an existing employee record by empId in MongoDB.
func (r *EmployeeRepository) Update(ctx context.Context, sub string, role string, empID string, emp model.Employee) (*model.Employee, error) {
	filter := bson.M{"empId": empID}
	if role != "admin" {
		filter["createdBy"] = sub
	}

	update := bson.M{
		"$set": bson.M{
			"fullName":   emp.FullName,
			"email":      emp.Email,
			"jobTitle":   emp.JobTitle,
			"department": emp.Department,
			"salary":     emp.Salary,
		},
		"$currentDate": bson.M{
			"updatedAt": true,
		},
	}

	var updatedEmp model.Employee
	err := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedEmp)

	if err != nil {
		return nil, err
	}

	return &updatedEmp, nil
}

// Delete removes an employee record by empId from MongoDB.
func (r *EmployeeRepository) Delete(ctx context.Context, sub string, role string, empID string) (*model.Employee, error) {
	filter := bson.M{"empId": empID}
	if role != "admin" {
		filter["createdBy"] = sub
	}

	var deletedEmp model.Employee
	err := r.collection.FindOneAndDelete(
		ctx,
		filter,
	).Decode(&deletedEmp)

	if err != nil {
		return nil, err
	}

	return &deletedEmp, nil
}

// Search queries employee records matching a regex across multiple fields.
func (r *EmployeeRepository) Search(ctx context.Context, sub string, role string, query string) ([]model.Employee, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"empId": bson.M{"$regex": query, "$options": "i"}},
			{"fullName": bson.M{"$regex": query, "$options": "i"}},
			{"email": bson.M{"$regex": query, "$options": "i"}},
			{"jobTitle": bson.M{"$regex": query, "$options": "i"}},
			{"department": bson.M{"$regex": query, "$options": "i"}},
		},
	}
	if role != "admin" {
		filter["createdBy"] = sub
	}

	// Sort by empId in ascending order
	findOptions := options.Find().SetSort(bson.D{{Key: "empId", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var employees []model.Employee
	if err := cursor.All(ctx, &employees); err != nil {
		return nil, err
	}

	if employees == nil {
		employees = []model.Employee{}
	}

	return employees, nil
}
