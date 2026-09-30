package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Employee struct {
	ID         bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EmpID      string        `json:"empId" bson:"empId"`
	FullName   string        `json:"fullName" bson:"fullName"`
	Email      string        `json:"email" bson:"email"`
	JobTitle   string        `json:"jobTitle" bson:"jobTitle"`
	Department string        `json:"department" bson:"department"`
	Salary     int           `json:"salary" bson:"salary"`
	CreatedBy  string        `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	CreatedAt  time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt" bson:"updatedAt"`
}
