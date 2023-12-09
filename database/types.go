package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Step represents a step in the process.
type Step struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Parameters []Parameter        `bson:"parameters,omitempty"`
	Function   string             `bson:"function,omitempty"`
	FunctionID string             `bson:"function_id,omitempty"`
}

// Parameter represents a parameter in a step.
type Parameter struct {
	Name  string      `bson:"name,omitempty"`
	Value interface{} `bson:"value,omitempty"`
	Type  string      `bson:"type,omitempty"`
}

// Trigger represents the trigger information.
type Trigger struct {
	Name       string             `bson:"name,omitempty"`
	Slug       string             `bson:"slug,omitempty"`
	Parameters []Parameter        `bson:"parameters,omitempty"`
	TriggerID  string             `bson:"trigger_id,omitempty"`
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt  time.Time          `bson:"created_at,omitempty"`
	UpdatedAt  time.Time          `bson:"updated_at,omitempty"`
}

// Workflow represents the overall workflow structure.
type Workflow struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty"`
	Steps     []Step             `bson:"steps,omitempty"`
	Trigger   Trigger            `bson:"trigger,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
	V         int                `bson:"__v,omitempty"`
}

type FlowRun struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FlowID    primitive.ObjectID `bson:"flow_id,omitempty"`
	Trigger   Trigger            `bson:"trigger,omitempty"`
	Steps     []StepRun          `bson:"steps,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
	V         int                `bson:"__v,omitempty"`
}

type StepRun struct {
	Name       string      `bson:"name,omitempty"`
	Parameters []Parameter `bson:"parameters,omitempty"`
	Function   string      `bson:"function,omitempty"`
	Logs       string      `bson:"logs,omitempty"`
	Success    bool        `bson:"success"`
	Message    string      `bson:"message,omitempty"`
	Value      interface{} `bson:"value,omitempty"`
	ID  primitive.ObjectID `bson:"_id,omitempty"`
}
