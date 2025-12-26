package main

import (
	"time"
)

// Service represents a monitored service
type Service struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Interval    int       `json:"interval"` // in seconds
	Status      string    `json:"status"`   // "online", "offline", "unknown"
	LastChecked time.Time `json:"lastChecked"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	// Downtime tracking
	WentOfflineAt  *time.Time `json:"wentOfflineAt,omitempty"`  // When service first went offline
	LastReminderAt *time.Time `json:"lastReminderAt,omitempty"` // When last reminder was sent
	// Failure tracking
	ConsecutiveFailures int `json:"consecutiveFailures"` // Number of consecutive failed checks
}

// ServiceStatus represents the current status of a service
type ServiceStatus struct {
	ServiceID    string    `json:"serviceId"`
	Status       string    `json:"status"`
	LastChecked  time.Time `json:"lastChecked"`
	ResponseTime int64     `json:"responseTime"` // in milliseconds
	Error        string    `json:"error,omitempty"`
}

// CheckResult represents the result of a health check
type CheckResult struct {
	ServiceID    string
	Status       string
	ResponseTime int64
	Error        error
	Timestamp    time.Time
}

// NotificationConfig holds Pushover configuration
type NotificationConfig struct {
	UserKey  string `json:"userKey"`
	AppToken string `json:"appToken"`
	Enabled  bool   `json:"enabled"`
}

// CreateServiceRequest represents the request to create a new service
type CreateServiceRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	URL      string `json:"url" validate:"required,httpurl,max=2048"`
	Interval int    `json:"interval" validate:"required,min=30,max=3600"`
}

// UpdateServiceRequest represents the request to update a service
type UpdateServiceRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	URL      string `json:"url" validate:"required,httpurl,max=2048"`
	Interval int    `json:"interval" validate:"required,min=30,max=3600"`
}

// BulkCreateServiceRequest represents a request to create multiple services
type BulkCreateServiceRequest struct {
	Services []CreateServiceRequest `json:"services" validate:"required,min=1,max=100,dive"`
}

// BulkUpdateServiceItem represents a single service update in a bulk request
type BulkUpdateServiceItem struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
	URL      string `json:"url" validate:"required,httpurl,max=2048"`
	Interval int    `json:"interval" validate:"required,min=30,max=3600"`
}

// BulkUpdateServiceRequest represents a request to update multiple services
type BulkUpdateServiceRequest struct {
	Services []BulkUpdateServiceItem `json:"services" validate:"required,min=1,max=100,dive"`
}

// BulkDeleteServiceRequest represents a request to delete multiple services
type BulkDeleteServiceRequest struct {
	IDs []string `json:"ids" validate:"required,min=1,max=100"`
}

// BulkOperationResponse represents the response for bulk operations
type BulkOperationResponse struct {
	Success  bool       `json:"success"`
	Count    int        `json:"count"`
	Services []*Service `json:"services,omitempty"`
}
