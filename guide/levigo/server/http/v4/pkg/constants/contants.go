package constants

import "time"

const (
	HttpPort   = "HTTP_PORT"
	ConfigPath = "CONFIG_PATH"
	RedisAddr  = "REDIS_ADDR"

	Yaml = "yaml"
	Tcp  = "tcp"

	STATUS   = "STATUS"
	HTTP     = "HTTP"
	ERROR    = "ERROR"
	METHOD   = "METHOD"
	METADATA = "METADATA"
	REQUEST  = "REQUEST"

	Page   = "page"
	Size   = "size"
	Search = "search"
	ID     = "id"

	PaymentsCollection      = "payments"
	HealthCheckKeyDefault   = "health-check:default"
	HealthCheckKeyFallback  = "health-check:fallback"
	HealthCheckTicker       = 1 * time.Second
	PaymentServiceHealthUrl = "/payments/service-health"
)
