package config

import (
	"flag"
	"os"
)

type Config struct {
	GCPProjectID string
	ServerAddr   string
}

func Load() *Config {
	var gcpProjectID string
	flag.StringVar(&gcpProjectID, "gcp-project-id", getEnv("GCP_PROJECT_ID", "demo-test"), "GCP Project ID")
	flag.Parse()

	return &Config{
		GCPProjectID: gcpProjectID,
		ServerAddr:   getEnv("SERVER_ADDR", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
