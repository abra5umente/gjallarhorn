package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// MonitorService handles service monitoring operations
type MonitorService struct {
	services map[string]*Service
	mu       sync.RWMutex
	client   *http.Client
	storage  *StorageService
}

// NewMonitorService creates a new monitor service
func NewMonitorService() *MonitorService {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // For self-signed certificates
			},
		},
	}

	storage := NewStorageService()

	// Load services from storage
	services, err := storage.LoadServices()
	if err != nil {
		log.Printf("Warning: Failed to load services from storage: %v", err)
		services = make(map[string]*Service)
	}

	return &MonitorService{
		services: services,
		client:   client,
		storage:  storage,
	}
}

// saveServices saves services to persistent storage
func (m *MonitorService) saveServices() error {
	m.mu.RLock()
	services := make(map[string]*Service)
	for k, v := range m.services {
		services[k] = v
	}
	m.mu.RUnlock()

	return m.storage.SaveServices(services)
}

// GetServices returns all monitored services
func (m *MonitorService) GetServices(c echo.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	services := make([]*Service, 0, len(m.services))
	for _, service := range m.services {
		services = append(services, service)
	}

	return c.JSON(http.StatusOK, services)
}

// CreateService creates a new monitored service
func (m *MonitorService) CreateService(c echo.Context) error {
	var req CreateServiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format: " + err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed: " + err.Error()})
	}

	service := &Service{
		ID:          uuid.New().String(),
		Name:        req.Name,
		URL:         req.URL,
		Interval:    req.Interval,
		Status:      "unknown",
		LastChecked: time.Time{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.mu.Lock()
	m.services[service.ID] = service
	m.mu.Unlock()

	// Save to persistent storage
	if err := m.saveServices(); err != nil {
		log.Printf("Warning: Failed to save services to storage: %v", err)
	}

	return c.JSON(http.StatusCreated, service)
}

// UpdateService updates an existing service
func (m *MonitorService) UpdateService(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "service ID is required"})
	}

	var req UpdateServiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format: " + err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed: " + err.Error()})
	}

	m.mu.Lock()
	service, exists := m.services[id]
	if !exists {
		m.mu.Unlock()
		return c.JSON(http.StatusNotFound, map[string]string{"error": "service not found"})
	}

	service.Name = req.Name
	service.URL = req.URL
	service.Interval = req.Interval
	service.UpdatedAt = time.Now()
	m.mu.Unlock()

	// Save to persistent storage
	if err := m.saveServices(); err != nil {
		log.Printf("Warning: Failed to save services to storage: %v", err)
	}

	return c.JSON(http.StatusOK, service)
}

// DeleteService deletes a service
func (m *MonitorService) DeleteService(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "service ID is required"})
	}

	m.mu.Lock()
	if _, exists := m.services[id]; !exists {
		m.mu.Unlock()
		return c.JSON(http.StatusNotFound, map[string]string{"error": "service not found"})
	}

	delete(m.services, id)
	m.mu.Unlock()

	// Save to persistent storage
	if err := m.saveServices(); err != nil {
		log.Printf("Warning: Failed to save services to storage: %v", err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetServiceStatus returns the current status of a service
func (m *MonitorService) GetServiceStatus(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "service ID is required"})
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	service, exists := m.services[id]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "service not found"})
	}

	status := ServiceStatus{
		ServiceID:   service.ID,
		Status:      service.Status,
		LastChecked: service.LastChecked,
	}

	return c.JSON(http.StatusOK, status)
}

// StartMonitoring starts the background monitoring process
func (m *MonitorService) StartMonitoring(notificationService *NotificationService) {
	// Service health check ticker (every 10 seconds)
	healthTicker := time.NewTicker(10 * time.Second)
	defer healthTicker.Stop()

	// Reminder check ticker (every hour)
	reminderTicker := time.NewTicker(1 * time.Hour)
	defer reminderTicker.Stop()

	for {
		select {
		case <-healthTicker.C:
			m.checkAllServices(notificationService)
		case <-reminderTicker.C:
			m.checkReminders(notificationService)
		}
	}
}

// checkAllServices checks all services for their health
func (m *MonitorService) checkAllServices(notificationService *NotificationService) {
	m.mu.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, service := range m.services {
		services = append(services, service)
	}
	m.mu.RUnlock()

	for _, service := range services {
		go m.checkService(service, notificationService)
	}
}

// checkService performs a health check on a single service
func (m *MonitorService) checkService(service *Service, notificationService *NotificationService) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", service.URL, nil)
	if err != nil {
		log.Printf("Error creating request for %s (%s): %v", service.Name, service.URL, err)
		m.updateServiceStatus(service, "failed", 0, err, notificationService)
		return
	}

	req.Header.Set("User-Agent", "Gjallarhorn/1.0")

	// Special handling for Plex
	if strings.Contains(strings.ToLower(service.URL), "plex") {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-Plex-Client-Identifier", "gjallarhorn-monitor")
	}

	resp, err := m.client.Do(req)
	responseTime := time.Since(start).Milliseconds()

	if err != nil {
		log.Printf("Request failed for %s (%s): %v", service.Name, service.URL, err)
		m.updateServiceStatus(service, "failed", responseTime, err, notificationService)
		return
	}
	defer resp.Body.Close()

	// Log the response details for debugging
	log.Printf("Service %s (%s): HTTP %d, Response time: %dms", service.Name, service.URL, resp.StatusCode, responseTime)

	// Consider 2xx, 3xx, and 401 (unauthorized) as healthy
	// 401 means the service is online but requires authentication
	if (resp.StatusCode >= 200 && resp.StatusCode < 400) || resp.StatusCode == 401 {
		if resp.StatusCode == 401 {
			log.Printf("Service %s (%s): HTTP 401 (unauthorized) - marking as online", service.Name, service.URL)
		}
		m.updateServiceStatus(service, "online", responseTime, nil, notificationService)
	} else {
		err := fmt.Errorf("HTTP %d", resp.StatusCode)
		log.Printf("Service %s (%s) failed with HTTP %d", service.Name, service.URL, resp.StatusCode)
		m.updateServiceStatus(service, "failed", responseTime, err, notificationService)
	}
}

// updateServiceStatus updates the service status and sends notifications if needed
func (m *MonitorService) updateServiceStatus(service *Service, status string, responseTime int64, err error, notificationService *NotificationService) {
	m.mu.Lock()
	defer m.mu.Unlock()

	previousStatus := service.Status
	service.LastChecked = time.Now()

	// Handle different status updates
	if status == "online" {
		// Service is healthy - reset failure counter and update status
		service.ConsecutiveFailures = 0

		if previousStatus == "offline" {
			// Service came back online - send recovery notification and clear downtime tracking
			var downtimeDuration string
			if service.WentOfflineAt != nil {
				duration := time.Since(*service.WentOfflineAt)
				if duration.Hours() >= 24 {
					days := int(duration.Hours() / 24)
					downtimeDuration = fmt.Sprintf("%d day(s)", days)
				} else if duration.Hours() >= 1 {
					hours := int(duration.Hours())
					downtimeDuration = fmt.Sprintf("%d hour(s)", hours)
				} else {
					minutes := int(duration.Minutes())
					downtimeDuration = fmt.Sprintf("%d minute(s)", minutes)
				}
			}

			// Send recovery notification
			notificationService.SendRecoveryNotification(service, downtimeDuration)

			// Clear downtime tracking
			service.WentOfflineAt = nil
			service.LastReminderAt = nil
		}

		service.Status = "online"

	} else if status == "failed" {
		// Service check failed - increment failure counter
		service.ConsecutiveFailures++
		log.Printf("Service %s (%s): Consecutive failures: %d/3", service.Name, service.URL, service.ConsecutiveFailures)

		// Only mark as offline after 3 consecutive failures
		if service.ConsecutiveFailures >= 3 {
			if previousStatus != "offline" {
				// Service just went offline - record the time and send initial notification
				now := time.Now()
				service.WentOfflineAt = &now
				service.LastReminderAt = &now // Set initial reminder time

				errorMsg := ""
				if err != nil {
					errorMsg = err.Error()
				}
				notificationService.SendNotification(service, errorMsg)
			}
			service.Status = "offline"
		} else {
			// Still within failure threshold - keep as online but log the failure
			service.Status = "online"
		}
	}
}

// checkReminders checks for services that have been down for over an hour and sends reminder notifications
func (m *MonitorService) checkReminders(notificationService *NotificationService) {
	m.mu.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, service := range m.services {
		services = append(services, service)
	}
	m.mu.RUnlock()

	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)

	for _, service := range services {
		// Only check services that are currently offline
		if service.Status != "offline" {
			continue
		}

		// Check if service has been down for over an hour
		if service.WentOfflineAt != nil && service.WentOfflineAt.Before(oneHourAgo) {
			// Check if we haven't sent a reminder in the last hour
			if service.LastReminderAt == nil || service.LastReminderAt.Before(oneHourAgo) {
				m.sendReminderNotification(service, notificationService)
			}
		}
	}
}

// sendReminderNotification sends a reminder notification for a service that's been down
func (m *MonitorService) sendReminderNotification(service *Service, notificationService *NotificationService) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Update the last reminder time
	now := time.Now()
	service.LastReminderAt = &now

	// Calculate downtime duration
	var downtimeDuration string
	if service.WentOfflineAt != nil {
		duration := now.Sub(*service.WentOfflineAt)
		if duration.Hours() >= 24 {
			days := int(duration.Hours() / 24)
			downtimeDuration = fmt.Sprintf("%d day(s)", days)
		} else if duration.Hours() >= 1 {
			hours := int(duration.Hours())
			downtimeDuration = fmt.Sprintf("%d hour(s)", hours)
		} else {
			minutes := int(duration.Minutes())
			downtimeDuration = fmt.Sprintf("%d minute(s)", minutes)
		}
	}

	// Send reminder notification
	notificationService.SendReminderNotification(service, downtimeDuration)
}
