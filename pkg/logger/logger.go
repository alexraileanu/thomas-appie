package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"

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
	l := Log{
		Message: message,
		Level:   "INFO",
		Extra:   extra,
	}

	s.send(l)
}

func (s *Service) Error(message string, error map[string]interface{}) {
	l := Log{
		Message: message,
		Level:   "ERROR",
		Extra:   error,
	}

	s.send(l)
}

func (s *Service) send(l Log) {
	if !s.enabled {
		logger := log.New(os.Stdout)
		if l.Level == "INFO" {
			logger.SetLevel(log.InfoLevel)
		} else if l.Level == "ERROR" {
			logger.SetLevel(log.ErrorLevel)
		}

		logger.Log(logger.GetLevel(), l.Message, l.Extra)
		return
	}
	body, err := json.Marshal(l)
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
