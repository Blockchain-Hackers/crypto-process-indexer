package main

import (
	"time"
)

// Step represents a step in the process.
type Step struct {
	Name       string      `json:"name"`
	Parameters []Parameter `json:"parameters"`
	Function   string      `json:"function"`
	FunctionID string      `json:"function_id"`
}

// Parameter represents a parameter in a step.
type Parameter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

// Trigger represents the trigger information.
type Trigger struct {
	Name       string      `json:"name"`
	Slug       string      `json:"slug"`
	Parameters []Parameter `json:"parameters"`
	TriggerID  string      `json:"trigger_id"`
	ID         string      `json:"_id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// Workflow represents the overall workflow structure.
type Workflow struct {
	ID struct {
		Oid string `json:"$oid"`
	} `json:"_id"`
	Name   string `json:"name"`
	UserID struct {
		Oid string `json:"$oid"`
	} `json:"user_id"`
	Steps     []Step    `json:"steps"`
	Trigger   Trigger   `json:"trigger"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	V         int       `json:"__v"`
}
