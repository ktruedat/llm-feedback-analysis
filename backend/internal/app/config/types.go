package config

import (
	"fmt"
	"strings"
)

// Profile represents the application running profile.
type Profile string

const (
	Development Profile = "dev"
	Production  Profile = "prod"
)

func (s Profile) Validate() error {
	switch s {
	case Development, Production:
		return nil
	default:
		return fmt.Errorf("invalid profile: %s", s)
	}
}

func (s Profile) String() string {
	return string(s)
}

type Server struct {
	Host                    string `yaml:"host" env:"HOST"`
	Port                    int    `yaml:"port" env:"PORT"`
	GracefulShutdownSeconds int    `yaml:"graceful_shutdown_seconds" env:"GRACEFUL_SHUTDOWN_SECONDS"`
}

func (s Server) Validate() error {
	if s.Port <= 0 || s.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", s.Port)
	}

	if s.GracefulShutdownSeconds < 0 {
		return fmt.Errorf("graceful shutdown seconds cannot be negative")
	}

	if strings.TrimSpace(s.Host) == "" {
		return fmt.Errorf("server host cannot be empty")
	}

	return nil
}

type Pagination struct {
	Limit  int `yaml:"limit" env:"LIMIT"`
	Offset int `yaml:"offset" env:"OFFSET"`
}

func (p Pagination) Validate() error {
	if p.Limit < 0 {
		return fmt.Errorf("pagination limit cannot be negative")
	}

	if p.Offset < 0 {
		return fmt.Errorf("pagination offset cannot be negative")
	}

	return nil
}

type Database struct {
	// DSN of the database, e.g., "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	DSN string `yaml:"dsn" env:"DSN"`
}

func (d Database) Validate() error {
	if strings.TrimSpace(d.DSN) == "" {
		return fmt.Errorf("database DSN cannot be empty")
	}

	return nil
}

type Tracing struct {
	// Enabled determines whether tracing is enabled.
	Enabled bool `yaml:"enabled" env:"ENABLED"`
	// OTelEndpoint is the OpenTelemetry collector endpoint (e.g., http://localhost:4318).
	// Required if Enabled is true.
	OTelEndpoint string `yaml:"otel_endpoint" env:"OTEL_ENDPOINT"`
	// TempoEndpoint is the Grafana Tempo endpoint (e.g., http://localhost:3200).
	// Required if Enabled is true.
	TempoEndpoint string `yaml:"tempo_endpoint" env:"TEMPO_ENDPOINT"`
	// ServiceName is the service name for tracing.
	// Required if Enabled is true.
	ServiceName string `yaml:"service_name" env:"SERVICE_NAME"`
	// ServiceVersion is the service version for tracing.
	// Required if Enabled is true.
	ServiceVersion string `yaml:"service_version" env:"SERVICE_VERSION"`
	// Insecure determines whether to use insecure connection (HTTP instead of HTTPS).
	// If using local OTel collector without TLS, set this to true.
	Insecure bool `yaml:"insecure" env:"INSECURE"`
}

func (t Tracing) Validate() error {
	if !t.Enabled {
		return nil
	}

	if strings.TrimSpace(t.OTelEndpoint) == "" {
		return fmt.Errorf("tracing is enabled but otel endpoint is not set")
	}
	if strings.TrimSpace(t.TempoEndpoint) == "" {
		return fmt.Errorf("tracing is enabled but tempo endpoint is not set")
	}

	if strings.TrimSpace(t.ServiceName) == "" {
		return fmt.Errorf("tracing is enabled but service name is not set")
	}

	if strings.TrimSpace(t.ServiceVersion) == "" {
		return fmt.Errorf("tracing is enabled but service version is not set")
	}

	return nil
}
