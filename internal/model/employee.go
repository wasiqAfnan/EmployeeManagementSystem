package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Employee struct {
	ID         bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EmpID      string        `json:"empId" bson:"empId"`
	FullName   string        `json:"fullName" bson:"fullName"`
	JobTitle   string        `json:"jobTitle" bson:"jobTitle"`
	Department string        `json:"department" bson:"department"`
	Salary     int           `json:"salary" bson:"salary"`
}
