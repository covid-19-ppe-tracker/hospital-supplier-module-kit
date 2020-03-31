package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type requirementType int

const (
	//SerializerContextKey is used to set data to be serialized in request context
	Mask requirementType = iota
	Gloves
)

//Requirement in-mem represents a requirements collection document
type Requirement struct {
	ID primitive.ObjectID `json:"-" bson:"_id,omitempty"`
}
