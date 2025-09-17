package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// StorageService handles persistent storage of services and configuration
type StorageService struct {
	servicesFile string
	configFile   string
	mu           sync.RWMutex
}

// NewStorageService creates a new storage service
func NewStorageService() *StorageService {
	return &StorageService{
		servicesFile: "/data/services.json",
		configFile:   "/data/config.json",
	}
}

// ensureDataDir creates the data directory if it doesn't exist
func (s *StorageService) ensureDataDir() error {
	return os.MkdirAll("/data", 0755)
}

// SaveServices saves services to persistent storage
func (s *StorageService) SaveServices(services map[string]*Service) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.ensureDataDir(); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	data, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal services: %v", err)
	}

	if err := os.WriteFile(s.servicesFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write services file: %v", err)
	}

	return nil
}

// LoadServices loads services from persistent storage
func (s *StorageService) LoadServices() (map[string]*Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.ensureDataDir(); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(s.servicesFile); os.IsNotExist(err) {
		// File doesn't exist, return empty map
		return make(map[string]*Service), nil
	}

	data, err := os.ReadFile(s.servicesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read services file: %v", err)
	}

	var services map[string]*Service
	if err := json.Unmarshal(data, &services); err != nil {
		return nil, fmt.Errorf("failed to unmarshal services: %v", err)
	}

	// Ensure we have a valid map
	if services == nil {
		services = make(map[string]*Service)
	}

	return services, nil
}

// SaveNotificationConfig saves notification configuration to persistent storage
func (s *StorageService) SaveNotificationConfig(config *NotificationConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.ensureDataDir(); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(s.configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// LoadNotificationConfig loads notification configuration from persistent storage
func (s *StorageService) LoadNotificationConfig() (*NotificationConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.ensureDataDir(); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(s.configFile); os.IsNotExist(err) {
		// File doesn't exist, return default config
		return &NotificationConfig{
			UserKey:  "",
			AppToken: "",
			Enabled:  false,
		}, nil
	}

	data, err := os.ReadFile(s.configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config NotificationConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}
