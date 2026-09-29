package repository

import (
	"context"

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

// GetAll fetches all employee records from MongoDB.
func (r *EmployeeRepository) GetAll(ctx context.Context) ([]model.Employee, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
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
func (r *EmployeeRepository) GetByEmpID(ctx context.Context, empID string) (*model.Employee, error) {
	var employee model.Employee
	err := r.collection.FindOne(ctx, bson.M{"empId": empID}).Decode(&employee)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// Create inserts a new employee record into MongoDB.
func (r *EmployeeRepository) Create(ctx context.Context, emp *model.Employee) (*model.Employee, error) {
	result, err := r.collection.InsertOne(ctx, emp)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		emp.ID = oid
	}
	return emp, nil
}

// Update updates an existing employee record by empId in MongoDB.
func (r *EmployeeRepository) Update(ctx context.Context, empID string, emp model.Employee) (*model.Employee, error) {
	filter := bson.M{"empId": empID}
	update := bson.M{
		"$set": bson.M{
			"fullName":   emp.FullName,
			"jobTitle":   emp.JobTitle,
			"department": emp.Department,
			"salary":     emp.Salary,
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
func (r *EmployeeRepository) Delete(ctx context.Context, empID string) (*model.Employee, error) {
	var deletedEmp model.Employee
	err := r.collection.FindOneAndDelete(
		ctx,
		bson.M{"empId": empID},
	).Decode(&deletedEmp)

	if err != nil {
		return nil, err
	}

	return &deletedEmp, nil
}

// Search queries employee records matching a regex across multiple fields.
func (r *EmployeeRepository) Search(ctx context.Context, query string) ([]model.Employee, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"empId": bson.M{"$regex": query, "$options": "i"}},
			{"fullName": bson.M{"$regex": query, "$options": "i"}},
			{"jobTitle": bson.M{"$regex": query, "$options": "i"}},
			{"department": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
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
