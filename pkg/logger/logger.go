package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

type Service struct {
	logger *log.Logger
}

type httpWriter struct {
	endpoint string
}

func (w *httpWriter) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest("POST", w.endpoint, bytes.NewBuffer(p))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-app-name", "thomas.appie")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending log: %v\n", err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error sending log: %v\n", resp.Status)
		return 0, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return len(p), nil
}

func New(enabled bool, debug bool) *Service {
	var output io.Writer
	if enabled {
		output = &httpWriter{endpoint: os.Getenv("LOGS_HOST")}
	} else {
		output = os.Stdout
	}

	l := log.New(output)
	l.SetFormatter(log.JSONFormatter)
	if debug {
		l.SetLevel(log.DebugLevel)
	} else {
		l.SetLevel(log.InfoLevel)
	}

	return &Service{logger: l}
}

func (s *Service) Info(message string, extra map[string]interface{}) {
	s.logger.Info(message, toArgs(extra)...)
}

func (s *Service) Debug(message string, extra map[string]interface{}) {
	s.logger.Debug(message, toArgs(extra)...)
}

func (s *Service) Error(message string, extra map[string]interface{}) {
	s.logger.Error(message, toArgs(extra)...)
}

func toArgs(extra map[string]interface{}) []interface{} {
	args := make([]interface{}, 0, len(extra)*2)
	for k, v := range extra {
		args = append(args, k, v)
	}
	return args
}
