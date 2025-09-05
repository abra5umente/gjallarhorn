package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

// NotificationService handles Pushover notifications
type NotificationService struct {
	config *NotificationConfig
}

// NewNotificationService creates a new notification service
func NewNotificationService() *NotificationService {
	config := &NotificationConfig{
		UserKey:  os.Getenv("PUSHOVER_USER_KEY"),
		AppToken: os.Getenv("PUSHOVER_APP_TOKEN"),
		Enabled:  os.Getenv("PUSHOVER_ENABLED") == "true",
	}

	return &NotificationService{
		config: config,
	}
}

// UpdateConfig updates the notification configuration
func (n *NotificationService) UpdateConfig(c echo.Context) error {
	var req NotificationConfig
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	n.config = &req
	return c.JSON(http.StatusOK, map[string]string{"message": "Notification configuration updated"})
}

// SendNotification sends a Pushover notification
func (n *NotificationService) SendNotification(service *Service, errorMsg string) {
	if !n.config.Enabled || n.config.UserKey == "" || n.config.AppToken == "" {
		return
	}

	title := fmt.Sprintf("🚨 Service Down: %s", service.Name)
	message := fmt.Sprintf("Service %s (%s) is currently offline.\nLast checked: %s",
		service.Name, service.URL, service.LastChecked.Format(time.RFC3339))

	if errorMsg != "" {
		message += fmt.Sprintf("\nError: %s", errorMsg)
	}

	payload := map[string]string{
		"token":    n.config.AppToken,
		"user":     n.config.UserKey,
		"title":    title,
		"message":  message,
		"sound":    "siren",
		"priority": "1", // High priority
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling notification payload: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post("https://api.pushover.net/1/messages.json",
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Pushover API error: %d\n", resp.StatusCode)
	}
}

// GetConfig returns the current notification configuration
func (n *NotificationService) GetConfig() *NotificationConfig {
	return n.config
}

// SendReminderNotification sends a reminder notification for a service that's been down
func (n *NotificationService) SendReminderNotification(service *Service, downtimeDuration string) {
	if !n.config.Enabled || n.config.UserKey == "" || n.config.AppToken == "" {
		return
	}

	title := fmt.Sprintf("⏰ Service Still Down: %s", service.Name)
	message := fmt.Sprintf("Service %s (%s) has been offline for %s.\nLast checked: %s\n\nThis is a reminder notification.",
		service.Name, service.URL, downtimeDuration, service.LastChecked.Format(time.RFC3339))

	payload := map[string]string{
		"token":    n.config.AppToken,
		"user":     n.config.UserKey,
		"title":    title,
		"message":  message,
		"sound":    "pushover", // Different sound for reminders
		"priority": "0",        // Normal priority for reminders
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling reminder notification payload: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post("https://api.pushover.net/1/messages.json",
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending reminder notification: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Pushover API error for reminder: %d\n", resp.StatusCode)
	}
}

// SendRecoveryNotification sends a notification when a service comes back online
func (n *NotificationService) SendRecoveryNotification(service *Service, downtimeDuration string) {
	if !n.config.Enabled || n.config.UserKey == "" || n.config.AppToken == "" {
		return
	}

	title := fmt.Sprintf("✅ Service Recovered: %s", service.Name)
	message := fmt.Sprintf("Service %s (%s) is back online!\nLast checked: %s",
		service.Name, service.URL, service.LastChecked.Format(time.RFC3339))

	if downtimeDuration != "" {
		message += fmt.Sprintf("\n\nTotal downtime: %s", downtimeDuration)
	}

	payload := map[string]string{
		"token":    n.config.AppToken,
		"user":     n.config.UserKey,
		"title":    title,
		"message":  message,
		"sound":    "magic", // Different sound for recovery
		"priority": "0",     // Normal priority
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling recovery notification payload: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post("https://api.pushover.net/1/messages.json",
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending recovery notification: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Pushover API error for recovery: %d\n", resp.StatusCode)
	}
}
