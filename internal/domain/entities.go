// Package domain defines entities for bff-api-go-collections-client-late-paying.
package domain

import "time"

// Activity represents a collection activity (gestión).
type Activity struct {
	ID           string    `json:"id"`
	AccountID    string    `json:"account_id"`
	ActivityType string    `json:"activity_type"`
	Result       string    `json:"result,omitempty"`
	Notes        string    `json:"notes"`
	PromiseDate  *string   `json:"promise_date,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Escalation represents an escalation request.
type Escalation struct {
	ID         string    `json:"id"`
	AccountID  string    `json:"account_id"`
	Reason     string    `json:"reason"`
	Status     string    `json:"status"` // pending, approved, rejected, assigned
	AssignedTo *string   `json:"assigned_to,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Account represents a portfolio account.
type Account struct {
	ID       string  `json:"id"`
	ClientID string  `json:"client_id"`
	LoanID   string  `json:"loan_id"`
	AgentID  string  `json:"agent_id"`
	Status   string  `json:"status"`
	Balance  float64 `json:"balance"`
	DPD      int     `json:"dpd"` // Days Past Due
}

// AgentPerformance represents agent performance metrics.
type AgentPerformance struct {
	AgentID          string  `json:"agent_id"`
	AgentName        string  `json:"agent_name"`
	ActivitiesCount  int     `json:"activities_count"`
	EscalationsCount int     `json:"escalations_count"`
	AssignmentsCount int     `json:"assignments_count"`
	PromiseCount     int     `json:"promise_count"`
	AvgResponseTime  float64 `json:"avg_response_time"`
}

// User represents an application user.
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"` // agent, supervisor, manager, admin, auditor
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
}

// Tenant represents a tenant (marca operativa).
type Tenant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Config represents system configuration.
type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
