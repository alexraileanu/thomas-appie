package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Service struct {
	enabled bool
}

type Log struct {
	Message string                 `json:"message"`
	Level   string                 `json:"level"`
	Extra   map[string]interface{} `json:"extra"`
}

func New(enabled bool) *Service {
	return &Service{
		enabled: enabled,
	}
}

func (s *Service) Info(message string, extra map[string]interface{}) {
	log := Log{
		Message: message,
		Level:   "INFO",
		Extra:   extra,
	}

	s.send(log)
}

func (s *Service) Error(message string, error map[string]interface{}) {
	log := Log{
		Message: message,
		Level:   "ERROR",
		Extra:   error,
	}

	s.send(log)
}

func (s *Service) send(log Log) {
	if !s.enabled {
		return
	}
	body, err := json.Marshal(log)
	if err != nil {
		fmt.Printf("Error marshalling log: %v", err)
		return
	}

	req, err := http.NewRequest("POST", os.Getenv("LOGS_HOST"), bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-app-name", "thomas.appie")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error sending log: %v", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error sending log: %v", resp.Status)
		return
	}
}
